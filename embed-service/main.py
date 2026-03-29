import base64
import io
import os
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer
from PIL import Image

# Reads SENTENCE_TRANSFORMERS_HOME from env — points to the baked-in model
MODEL_NAME = "clip-ViT-B-32"
model = SentenceTransformer(MODEL_NAME)  # loads from /app/model-store, no download needed

app = FastAPI()

class EmbedRequest(BaseModel):
    text: str | None = None
    image_b64: str | None = None

@app.post("/embed")
async def embed(req: EmbedRequest):
    try:
        if req.image_b64:
            img_data = base64.b64decode(req.image_b64)
            input_data = Image.open(io.BytesIO(img_data)).convert("RGB")
        elif req.text:
            input_data = req.text
        else:
            raise HTTPException(status_code=400, detail="Provide 'text' or 'image_b64'")

        vector = model.encode(input_data, convert_to_numpy=True).tolist()
        return {"vector": vector, "dimensions": len(vector)}

    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
def health():
    return {"status": "ready", "model": MODEL_NAME, "dims": 512}