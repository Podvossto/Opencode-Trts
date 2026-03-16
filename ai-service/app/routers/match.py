# Purpose: Match router — compute resume-JD similarity score using hard skill weights
# Ref: API Contract — POST /match
# Hard skill weights: ratios 1-5 (not required to sum to 100)
# If no hard skills defined for job → JD similarity score = 100%
# Dev must implement: weighted cosine similarity between resume embedding and JD
from fastapi import APIRouter
from pydantic import BaseModel
from typing import List, Optional

router = APIRouter()


class HardSkillWeight(BaseModel):
    skill_name: str
    weight: int  # 1-5 ratio


class MatchRequest(BaseModel):
    application_id: str
    resume_embedding: List[float]
    jd_embedding: List[float]
    hard_skills: Optional[List[HardSkillWeight]] = None  # empty = 100% JD similarity


class MatchResponse(BaseModel):
    application_id: str
    similarity_score: float      # 0.0 - 1.0
    hard_skill_score: float      # 0.0 - 1.0
    combined_score: float        # weighted combination


@router.post("", response_model=MatchResponse)
async def compute_match(req: MatchRequest) -> MatchResponse:
    """
    Compute resume-JD match score.
    If hard_skills is empty/None: combined_score = JD cosine similarity (100% weight).
    If hard_skills present: combined_score = weighted sum of (JD similarity + hard skill ratios).
    Dev must implement: cosine similarity + hard skill weighted scoring logic.
    """
    # Dev must implement matching logic
    raise NotImplementedError("Match scoring not implemented")
