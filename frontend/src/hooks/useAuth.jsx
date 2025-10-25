import { useState, useEffect } from 'preact/hooks';
import apiClient from '../utils/apiClient';

export function useAuth() {
    const [user, setUser] = useState(null);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        // Проверяем наличие токена при загрузке
        const token = apiClient.getToken();
        if (token) {
            setIsAuthenticated(true);
            // TODO: можно добавить запрос /auth/me для получения данных пользователя
        }
        setIsLoading(false);
        
        // Глобальный обработчик ошибок 401
        const handleUnauthorized = () => {
            setIsAuthenticated(false);
            setUser(null);
        };
        
        // Добавляем обработчик для глобальных ошибок 401
        window.addEventListener('unauthorized', handleUnauthorized);
        
        return () => {
            window.removeEventListener('unauthorized', handleUnauthorized);
        };
    }, []);

    const login = async (email, password) => {
        try {
            const response = await apiClient.login(email, password);
            setUser(response.user || { email });
            setIsAuthenticated(true);
            return { success: true };
        } catch (error) {
            return { success: false, error: error.message };
        }
    };

    const register = async (email, password) => {
        try {
            await apiClient.register(email, password);
            // После регистрации автоматически логинимся
            return await login(email, password);
        } catch (error) {
            return { success: false, error: error.message };
        }
    };

    const logout = () => {
        apiClient.logout();
        setUser(null);
        setIsAuthenticated(false);
    };

    return {
        user,
        isAuthenticated,
        isLoading,
        login,
        register,
        logout,
    };
}

