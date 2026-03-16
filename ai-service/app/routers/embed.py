# Purpose: Embedding router — generate vector embeddings from extracted resume text
# Ref: API Contract — POST /embed
# Dev must implement: sentence-transformers encoding, store vector in DB
from fastapi import APIRouter
from pydantic import BaseModel
from typing import List

router = APIRouter()


class EmbedRequest(BaseModel):
    application_id: str
    text: str  # Extracted text from OCR


class EmbedResponse(BaseModel):
    application_id: str
    embedding: List[float]  # 384-dimensional vector (sentence-transformers default)
    model: str


@router.post("", response_model=EmbedResponse)
async def generate_embedding(req: EmbedRequest) -> EmbedResponse:
    """
    Generate a sentence embedding vector from resume text.
    Dev must implement: load sentence-transformers model, encode text,
    persist embedding to applications table.
    """
    # Dev must implement: model = SentenceTransformer('all-MiniLM-L6-v2')
    raise NotImplementedError("Embedding generation not implemented")
