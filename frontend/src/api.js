// Базовый URL для API
export const API_BASE_URL = 'http://109.73.202.151:8080' // Например, 'http://localhost:8080' или '' для относительных путей

// Функция для получения полного URL
export function apiUrl(path) {
  return `${API_BASE_URL}${path}`;
} 