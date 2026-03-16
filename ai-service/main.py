# Purpose: FastAPI AI service entry point — OCR and embedding endpoints
# Ref: API Contract — POST /ocr, POST /embed, POST /match, GET /health
# Dev must implement: OCR extraction (Tesseract/PyMuPDF), sentence embedding (sentence-transformers),
#                     JD similarity scoring, request validation
from fastapi import FastAPI
from app.routers import ocr, embed, match
from app.core.config import settings

app = FastAPI(
    title="ATS AI Service",
    description="OCR extraction, text embedding, and resume-JD matching service",
    version="0.1.0",
)

# Register routers
app.include_router(ocr.router, prefix="/ocr", tags=["OCR"])
app.include_router(embed.router, prefix="/embed", tags=["Embedding"])
app.include_router(match.router, prefix="/match", tags=["Matching"])


@app.get("/health")
async def health_check():
    return {"status": "ok"}
