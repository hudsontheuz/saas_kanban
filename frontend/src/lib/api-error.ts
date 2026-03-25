import axios from 'axios';

export function getErrorMessage(error: unknown, fallback = 'Ocorreu um erro inesperado.'): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data;
    if (data && typeof data === 'object') {
      const apiError = 'error' in data ? data.error : undefined;
      const apiMessage = 'message' in data ? data.message : undefined;
      if (typeof apiError === 'string' && apiError.trim()) return apiError;
      if (typeof apiMessage === 'string' && apiMessage.trim()) return apiMessage;
    }

    if (typeof error.message === 'string' && error.message.trim()) {
      return error.message;
    }
  }

  if (error instanceof Error && error.message.trim()) {
    return error.message;
  }

  return fallback;
}
