FROM python:3.11-slim

# Single source of truth for model location
ENV SENTENCE_TRANSFORMERS_HOME=/app/model-store

WORKDIR /app

# System deps for Pillow/OpenCV
RUN apt-get update && apt-get install -y \
    libgl1 \
    libglib2.0-0 \
    curl \
    && rm -rf /var/lib/apt/lists/*

# CPU PyTorch first — prevents pulling the 2GB GPU wheel
RUN pip install --no-cache-dir torch torchvision \
    --index-url https://download.pytorch.org/whl/cpu

# App dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Bake the model into the image at /app/model-store
# Volume is NOT mounted here, so this layer survives at runtime
RUN python -c "from sentence_transformers import SentenceTransformer; SentenceTransformer('clip-ViT-B-32')"

COPY main.py .

EXPOSE 8001

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8001", "--workers", "1"]