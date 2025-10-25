import { useState, useEffect } from 'preact/hooks';
import { motion, AnimatePresence } from 'framer-motion';
import apiClient from '../utils/apiClient';

export function Projects({ onSelectProject }) {
    const [projects, setProjects] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState('');
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [newProjectName, setNewProjectName] = useState('');
    const [isCreating, setIsCreating] = useState(false);

    useEffect(() => {
        loadProjects();
    }, []);

    const loadProjects = async () => {
        try {
            setIsLoading(true);
            const data = await apiClient.getProjects();
            setProjects(data || []);
        } catch (err) {
            // Не показываем ошибку 401 - она обрабатывается глобально
            if (err.message && !err.message.includes('HTTP 401')) {
                setError(err.message);
            }
        } finally {
            setIsLoading(false);
        }
    };

    const handleCreateProject = async (e) => {
        e.preventDefault();
        if (!newProjectName.trim()) return;

        try {
            setIsCreating(true);
            const project = await apiClient.createProject({ name: newProjectName });
            setProjects([...projects, project]);
            setNewProjectName('');
            setShowCreateModal(false);
        } catch (err) {
            // Не показываем ошибку 401 - она обрабатывается глобально
            if (err.message && !err.message.includes('HTTP 401')) {
                setError(err.message);
            }
        } finally {
            setIsCreating(false);
        }
    };

    const copyToClipboard = (text) => {
        navigator.clipboard.writeText(text);
        // TODO: показать toast уведомление
    };

    if (isLoading) {
        return (
            <div className="flex items-center justify-center py-12">
                <div className="text-gray-400">Загрузка проектов...</div>
            </div>
        );
    }

    return (
        <div>
            {/* Header */}
            <div className="flex justify-between items-center mb-6">
                <div>
                    <h2 className="text-2xl font-bold text-white mb-1">Ваши проекты</h2>
                    <p className="text-gray-400 text-sm">
                        Создавайте проекты и получайте API для mock данных
                    </p>
                </div>
                <motion.button
                    onClick={() => setShowCreateModal(true)}
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                    className="btn-primary"
                >
                    + Новый проект
                </motion.button>
            </div>

            {/* Error */}
            {error && (
                <motion.div
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="bg-red-900/20 border border-red-800 rounded-lg p-4 text-red-400 mb-6"
                >
                    {error}
                </motion.div>
            )}

            {/* Projects Grid */}
            {projects.length === 0 ? (
                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="card text-center py-12"
                >
                    <div className="text-gray-500 mb-4">
                        <svg className="w-16 h-16 mx-auto mb-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
                        </svg>
                    </div>
                    <h3 className="text-xl font-semibold text-white mb-2">
                        Нет проектов
                    </h3>
                    <p className="text-gray-400 mb-6">
                        Создайте первый проект для начала работы
                    </p>
                    <button
                        onClick={() => setShowCreateModal(true)}
                        className="btn-primary"
                    >
                        Создать проект
                    </button>
                </motion.div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {projects.map((project, index) => (
                        <motion.div
                            key={project.id}
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ delay: index * 0.05 }}
                            className="card hover:border-primary-600 transition-colors cursor-pointer"
                            onClick={() => onSelectProject(project)}
                        >
                            <h3 className="text-lg font-semibold text-white mb-3">
                                {project.name}
                            </h3>
                            
                            <div className="space-y-2 text-sm">
                                <div>
                                    <span className="text-gray-400 block mb-1">API URL:</span>
                                    <div className="flex items-center gap-2 bg-dark-900 rounded px-2 py-1.5">
                                        <code className="text-primary-400 font-mono text-xs flex-1 truncate">
                                            {project.base_url || `http://localhost:8080/${project.api_key}`}
                                        </code>
                                        <button
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                copyToClipboard(project.base_url || `http://localhost:8080/${project.api_key}`);
                                            }}
                                            className="text-gray-400 hover:text-white transition-colors flex-shrink-0"
                                            title="Копировать URL"
                                        >
                                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                                            </svg>
                                        </button>
                                    </div>
                                </div>
                                
                                <div className="flex items-center justify-between pt-2 border-t border-dark-600">
                                    <span className="text-gray-400">Коллекций:</span>
                                    <span className="text-white font-medium">
                                        {project.collections_count || 0}
                                    </span>
                                </div>
                            </div>
                        </motion.div>
                    ))}
                </div>
            )}

            {/* Create Modal */}
            <AnimatePresence>
                {showCreateModal && (
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50"
                        onClick={() => setShowCreateModal(false)}
                    >
                        <motion.div
                            initial={{ scale: 0.95, y: 20 }}
                            animate={{ scale: 1, y: 0 }}
                            exit={{ scale: 0.95, y: 20 }}
                            className="card max-w-md w-full"
                            onClick={(e) => e.stopPropagation()}
                        >
                            <h3 className="text-xl font-bold text-white mb-4">
                                Новый проект
                            </h3>
                            
                            <form onSubmit={handleCreateProject}>
                                <div className="mb-4">
                                    <label className="block text-sm font-medium text-gray-300 mb-2">
                                        Название проекта
                                    </label>
                                    <input
                                        type="text"
                                        value={newProjectName}
                                        onInput={(e) => setNewProjectName(e.target.value)}
                                        className="input"
                                        placeholder="My Project"
                                        required
                                        disabled={isCreating}
                                        autoFocus
                                    />
                                </div>

                                <div className="flex gap-3">
                                    <button
                                        type="button"
                                        onClick={() => setShowCreateModal(false)}
                                        className="btn-ghost flex-1"
                                        disabled={isCreating}
                                    >
                                        Отмена
                                    </button>
                                    <button
                                        type="submit"
                                        className="btn-primary flex-1"
                                        disabled={isCreating}
                                    >
                                        {isCreating ? 'Создание...' : 'Создать'}
                                    </button>
                                </div>
                            </form>
                        </motion.div>
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    );
}

