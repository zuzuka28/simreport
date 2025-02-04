from abc import ABC, abstractmethod

import numpy as np


class DocumentVectorizerService(ABC):
    @abstractmethod
    async def vectorize_text(self, text: str) -> np.ndarray:
        """vectorize text document"""
