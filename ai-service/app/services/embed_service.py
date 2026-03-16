# Purpose: Embedding service — generates vector embeddings for resume text and job descriptions
# Ref: ai-service/app/routers/embed.py — POST /embed/text
# Dev must implement:
#   - get_embedding(text: str) -> list[float]
#     Uses sentence-transformers (paraphrase-multilingual-MiniLM-L12-v2 for EN+TH support)
#     Returns 384-dimensional float vector
#   - get_embeddings_batch(texts: list[str]) -> list[list[float]]
#     Batch encode for efficiency — max 32 texts per call
#   - Model is loaded once at module level (singleton) to avoid reload overhead
# Dependencies: sentence-transformers, torch (in requirements.txt)

from typing import Optional

# TODO: Dev must load model at module level as singleton
# model = SentenceTransformer('paraphrase-multilingual-MiniLM-L12-v2')


def get_embedding(text: str) -> list:
    """
    Generate a single 384-dim embedding vector for the given text.
    TODO: Implement using sentence-transformers model.
    """
    raise NotImplementedError("get_embedding not yet implemented")


def get_embeddings_batch(texts: list) -> list:
    """
    Generate embeddings for a batch of texts (max 32).
    TODO: Implement using model.encode(texts, batch_size=32).
    Returns list of 384-dim float vectors.
    """
    raise NotImplementedError("get_embeddings_batch not yet implemented")
