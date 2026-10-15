from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    DEBUG: bool = True
    model_config = SettingsConfigDict(
        
        env_file='.env', env_file_encoding='utf-8'
    )
    DATABASE_URL: str
