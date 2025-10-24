const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

class APIClient {
    constructor() {
        this.baseURL = API_BASE_URL;
        this.token = localStorage.getItem('token');
    }

    setToken(token) {
        this.token = token;
        if (token) {
            localStorage.setItem('token', token);
        } else {
            localStorage.removeItem('token');
        }
    }

    getToken() {
        return this.token;
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;

        const headers = {
            'Content-Type': 'application/json',
            ...options.headers,
        };

        // Всегда читаем актуальный токен (на случай если обновился после логина)
        const currentToken = this.token || localStorage.getItem('token');
        if (currentToken) {
            headers['Authorization'] = `Bearer ${currentToken}`;
        }

        const config = {
            ...options,
            headers,
        };

        try {
            const response = await fetch(url, config);

            if (!response.ok) {
                const error = await response.json().catch(() => ({ message: response.statusText }));
                throw new Error(error.message || `HTTP ${response.status}`);
            }

            // Для 204 No Content не парсим JSON
            if (response.status === 204) {
                return null;
            }

            return await response.json();
        } catch (error) {
            console.error(`API Error [${endpoint}]:`, error);
            throw error;
        }
    }

    // Auth endpoints
    async register(email, password) {
        return this.request('/auth/register', {
            method: 'POST',
            body: JSON.stringify({ email, password }),
        });
    }

    async login(email, password) {
        const response = await this.request('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ email, password }),
        });

        // Backend возвращает access_token (snake_case)
        if (response.access_token) {
            this.setToken(response.access_token);
        }

        return response;
    }

    logout() {
        this.setToken(null);
    }

    // Projects endpoints
    async getProjects() {
        const response = await this.request('/projects');
        // Backend возвращает { projects: [...], count: N }
        return response.projects || [];
    }

    async getProject(id) {
        return this.request(`/projects/${id}`);
    }

    async createProject(data) {
        return this.request('/projects', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateProject(id, data) {
        return this.request(`/projects/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteProject(id) {
        return this.request(`/projects/${id}`, {
            method: 'DELETE',
        });
    }

    // Collections endpoints
    async getCollections(projectId) {
        const response = await this.request(`/projects/${projectId}/collections`);
        // Backend возвращает { collections: [...], count: N }
        return response.collections || [];
    }

    async createCollection(projectId, data) {
        return this.request(`/projects/${projectId}/collections`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateCollection(projectId, collectionId, data) {
        return this.request(`/projects/${projectId}/collections/${collectionId}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteCollection(projectId, collectionId) {
        return this.request(`/projects/${projectId}/collections/${collectionId}`, {
            method: 'DELETE',
        });
    }

    // Data endpoints (public API с api_key)
    async getCollectionData(apiKey, collectionName) {
        const response = await this.request(`/${apiKey}/${collectionName}`);
        // Backend возвращает чистый массив документов
        return response;
    }

    async createDocument(apiKey, collectionName, data) {
        return this.request(`/${apiKey}/${collectionName}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateDocument(apiKey, collectionName, documentId, data) {
        return this.request(`/${apiKey}/${collectionName}/${documentId}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteDocument(apiKey, collectionName, documentId) {
        return this.request(`/${apiKey}/${collectionName}/${documentId}`, {
            method: 'DELETE',
        });
    }

    async flushCollection(apiKey, collectionName) {
        return this.request(`/${apiKey}/${collectionName}`, {
            method: 'DELETE',
        });
    }
}

export default new APIClient();

