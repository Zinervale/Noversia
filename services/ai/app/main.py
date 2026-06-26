from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI(title="Noversia AI Service", version="0.1.0")

class AnalyzeRequest(BaseModel):
    message: str
    context: dict | None = None

class AnalyzeResponse(BaseModel):
    answer: str
    confidence_score: float
    source: str

@app.get("/health")
def health():
    return {"status": "ok", "service": "noversia-ai"}

@app.post("/analyze", response_model=AnalyzeResponse)
def analyze(request: AnalyzeRequest):
    return AnalyzeResponse(
        answer="Analyse IA simulée : les données fournies ne suffisent pas encore pour produire une recommandation définitive.",
        confidence_score=0.42,
        source="mock"
    )
