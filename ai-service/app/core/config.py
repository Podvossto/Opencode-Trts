# Purpose: AI service settings — loads environment variables
# Dev must implement: DATABASE_URL if direct DB write needed, AI_MODEL_NAME, port
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    ai_model_name: str = "all-MiniLM-L6-v2"
    go_api_url: str = "http://go-api:8080"
    port: int = 8000
    log_level: str = "info"

    class Config:
        env_file = ".env"


settings = Settings()
