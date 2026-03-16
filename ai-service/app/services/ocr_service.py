# Purpose: OCR service — extracts text from uploaded resume PDFs/images
# Ref: ai-service/app/routers/ocr.py — POST /ocr/extract
# Dev must implement:
#   - extract_text_from_file(file_bytes: bytes, mime_type: str) -> str
#     Uses pytesseract for images, pdfplumber for PDFs
#     Returns raw extracted text (UTF-8 string)
#   - preprocess_image(image: PIL.Image) -> PIL.Image
#     Greyscale + threshold for better OCR accuracy on low-quality scans
# Dependencies: pytesseract, Pillow, pdfplumber (all in requirements.txt)

from typing import Optional


def extract_text_from_file(file_bytes: bytes, mime_type: str) -> str:
    """
    Extract raw text from a PDF or image file.
    TODO: Implement using pdfplumber (PDF) or pytesseract (image).
    """
    raise NotImplementedError("extract_text_from_file not yet implemented")


def preprocess_image(image_bytes: bytes) -> bytes:
    """
    Preprocess image for improved OCR accuracy.
    TODO: Implement greyscale + threshold using Pillow.
    Returns preprocessed image as bytes.
    """
    raise NotImplementedError("preprocess_image not yet implemented")
