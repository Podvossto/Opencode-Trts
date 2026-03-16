# Purpose: Match service — computes cosine similarity score between resume and job description embeddings
# Ref: ai-service/app/routers/match.py — POST /match/score
# Dev must implement:
#   - cosine_similarity(vec_a: list[float], vec_b: list[float]) -> float
#     Pure numpy implementation — returns value in [0.0, 1.0]
#   - score_resume(resume_text: str, job_description: str) -> MatchResult
#     Calls embed_service.get_embedding() for both texts, then cosine_similarity
#     Returns MatchResult with score (0-100 int) and reasoning dict
#   - dataclass MatchResult { score: int, raw_similarity: float, reasoning: dict }
#     reasoning keys: keyword_overlap, semantic_score, experience_signal
# Dependencies: numpy, embed_service (internal)

from dataclasses import dataclass, field
from typing import Optional


@dataclass
class MatchResult:
    score: int                    # 0-100 integer match score
    raw_similarity: float         # cosine similarity 0.0-1.0
    reasoning: dict = field(default_factory=dict)
    # reasoning keys: keyword_overlap (float), semantic_score (float), experience_signal (str)


def cosine_similarity(vec_a: list, vec_b: list) -> float:
    """
    Compute cosine similarity between two vectors.
    TODO: Implement using numpy dot product / (norm_a * norm_b).
    Returns float in [0.0, 1.0].
    """
    raise NotImplementedError("cosine_similarity not yet implemented")


def score_resume(resume_text: str, job_description: str) -> MatchResult:
    """
    Score a resume against a job description.
    TODO: Implement — call embed_service.get_embedding() for both, compute cosine_similarity,
    scale to 0-100, populate reasoning dict.
    Returns MatchResult.
    """
    raise NotImplementedError("score_resume not yet implemented")
