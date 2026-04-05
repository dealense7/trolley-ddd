import base64
import io
import torch
import numpy as np
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer
from PIL import Image

MODEL_NAME = "clip-ViT-B-32"
TEXT_MODEL  = "all-MiniLM-L6-v2"

device = "cuda" if torch.cuda.is_available() else "cpu"
clip_model = SentenceTransformer(MODEL_NAME, device=device)
text_model = SentenceTransformer(TEXT_MODEL, device=device)

app = FastAPI()

class EmbedRequest(BaseModel):
    text: str | None = None
    image_b64: str | None = None
    text_weight: float = 0.55
    image_weight: float = 0.45

@app.post("/embed")
async def embed(req: EmbedRequest):
    # Bug 3 fix: validate that at least one input is provided
    if req.text is None and req.image_b64 is None:
        raise HTTPException(status_code=422, detail="At least one of 'text' or 'image_b64' must be provided.")

    try:
        if req.image_b64:
            img_data = base64.b64decode(req.image_b64)
            img = Image.open(io.BytesIO(img_data)).convert("RGB")
            img_vec = clip_model.encode(img, convert_to_numpy=True)

            if req.text:
                # Bug 1 + 2 fix: use clip_model for text when fusing with an image
                # so both vectors are in the same 512-dim CLIP space.
                text_vec = clip_model.encode(req.text, convert_to_numpy=True)

                text_norm = text_vec / (np.linalg.norm(text_vec) + 1e-9)
                img_norm  = img_vec  / (np.linalg.norm(img_vec)  + 1e-9)
                fused = req.text_weight * text_norm + req.image_weight * img_norm
            else:
                # Image-only: no text to fuse
                fused = img_vec
        else:
            # Text-only: MiniLM is fine here (no fusion, no dim conflict)
            fused = text_model.encode(req.text, convert_to_numpy=True)

        fused = fused / (np.linalg.norm(fused) + 1e-9)
        return {"vector": fused.tolist(), "dimensions": len(fused)}

    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
def health():
    return {
        "status": "ready",
        "clip_model": MODEL_NAME,
        "text_model": TEXT_MODEL,
        "clip_dims": 512,
        "text_dims": 384,
        "device": device
    }