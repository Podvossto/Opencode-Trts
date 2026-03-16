# Purpose: OCR router — extract text from uploaded resume documents
# Ref: API Contract — POST /ocr
# Dev must implement: extract text from PDF/image using PyMuPDF + Tesseract
from fastapi import APIRouter
from pydantic import BaseModel

router = APIRouter()


class OCRRequest(BaseModel):
    application_id: str
    document_url: str  # Pre-signed URL or file path


class OCRResponse(BaseModel):
    application_id: str
    extracted_text: str
    page_count: int


@router.post("", response_model=OCRResponse)
async def extract_text(req: OCRRequest) -> OCRResponse:
    """
    Extract text from a resume document (PDF or image).
    Dev must implement: download from document_url, extract with PyMuPDF/Tesseract,
    return extracted_text for embedding pipeline.
    """
    # Dev must implement extraction logic
    raise NotImplementedError("OCR extraction not implemented")
