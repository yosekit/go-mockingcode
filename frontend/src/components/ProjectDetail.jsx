import { useState } from 'preact/hooks';
import { motion, AnimatePresence } from 'framer-motion';
import apiClient from '../utils/apiClient';
import { Collections } from './Collections';

export function ProjectDetail({ project, onBack, onProjectUpdated, onProjectDeleted }) {
    const [isEditing, setIsEditing] = useState(false);
    const [projectName, setProjectName] = useState(project.name);
    const [projectDescription, setProjectDescription] = useState(project.description || '');
    const [isSaving, setIsSaving] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);
    const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
    const [error, setError] = useState('');

    const handleSave = async (e) => {
        e.preventDefault();
        if (!projectName.trim()) return;

        try {
            setIsSaving(true);
            setError('');
            const updated = await apiClient.updateProject(project.id, {
                name: projectName,
                description: projectDescription,
            });
            onProjectUpdated({ ...project, ...updated });
            setIsEditing(false);
        } catch (err) {
            // Не показываем ошибку 401 - она обрабатывается глобально
            if (err.message && !err.message.includes('HTTP 401')) {
                setError(err.message);
            }
        } finally {
            setIsSaving(false);
        }
    };

    const handleDelete = async () => {
        try {
            setIsDeleting(true);
            await apiClient.deleteProject(project.id);
            onProjectDeleted(project.id);
        } catch (err) {
            setError(err.message);
            setIsDeleting(false);
        }
    };

    const copyToClipboard = (text) => {
        navigator.clipboard.writeText(text);
        // TODO: показать toast уведомление
    };

    const apiUrl = project.base_url || `http://localhost:8080/${project.api_key}`;

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="card">
                <div className="flex items-start justify-between mb-4">
                    <div className="flex-1">
                        {isEditing ? (
                            <form onSubmit={handleSave} className="space-y-3">
                                <div>
                                    <label className="block text-sm font-medium text-gray-300 mb-2">
                                        Название проекта
                                    </label>
                                    <input
                                        type="text"
                                        value={projectName}
                                        onInput={(e) => setProjectName(e.target.value)}
                                        className="input"
                                        placeholder="My Project"
                                        required
                                        disabled={isSaving}
                                        autoFocus
                                    />
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-300 mb-2">
                                        Описание (опционально)
                                    </label>
                                    <textarea
                                        value={projectDescription}
                                        onInput={(e) => setProjectDescription(e.target.value)}
                                        className="input resize-none"
                                        placeholder="Краткое описание проекта"
                                        rows={3}
                                        disabled={isSaving}
                                    />
                                </div>

                                <div className="flex gap-3">
                                    <button
                                        type="button"
                                        onClick={() => {
                                            setIsEditing(false);
                                            setProjectName(project.name);
                                            setProjectDescription(project.description || '');
                                            setError('');
                                        }}
                                        className="btn-ghost"
                                        disabled={isSaving}
                                    >
                                        Отмена
                                    </button>
                                    <button
                                        type="submit"
                                        className="btn-primary"
                                        disabled={isSaving}
                                    >
                                        {isSaving ? 'Сохранение...' : 'Сохранить'}
                                    </button>
                                </div>
                            </form>
                        ) : (
                            <>
                                <div className="flex items-start justify-between mb-4">
                                    <div>
                                        <h2 className="text-2xl font-bold text-white mb-2">
                                            {project.name}
                                        </h2>
                                        {project.description && (
                                            <p className="text-gray-400 text-sm">
                                                {project.description}
                                            </p>
                                        )}
                                    </div>
                                    <button
                                        onClick={() => setIsEditing(true)}
                                        className="btn-ghost"
                                    >
                                        Редактировать
                                    </button>
                                </div>

                                {/* API URL */}
                                <div className="space-y-2">
                                    <span className="text-gray-400 text-sm block">API URL проекта:</span>
                                    <div className="flex items-center gap-2 bg-dark-900 rounded-lg px-3 py-2.5">
                                        <code className="text-primary-400 font-mono text-sm flex-1">
                                            {apiUrl}
                                        </code>
                                        <button
                                            onClick={() => copyToClipboard(apiUrl)}
                                            className="text-gray-400 hover:text-white transition-colors"
                                            title="Копировать URL"
                                        >
                                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                                            </svg>
                                        </button>
                                    </div>
                                    <p className="text-gray-500 text-xs">
                                        Используйте этот URL для доступа к API вашего проекта
                                    </p>
                                </div>
                            </>
                        )}
                    </div>
                </div>

                {/* Error */}
                {error && (
                    <motion.div
                        initial={{ opacity: 0, y: -10 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="bg-red-900/20 border border-red-800 rounded-lg p-3 text-red-400 text-sm"
                    >
                        {error}
                    </motion.div>
                )}
            </div>

            {/* Коллекции */}
            <div className="card">
                <Collections projectId={project.id} apiKey={project.api_key} />
            </div>

            {/* Danger Zone */}
            <div className="card border-red-900/50">
                <h3 className="text-lg font-semibold text-red-400 mb-2">
                    Опасная зона
                </h3>
                <p className="text-gray-400 text-sm mb-4">
                    Удаление проекта необратимо. Все данные будут потеряны.
                </p>
                <button
                    onClick={() => setShowDeleteConfirm(true)}
                    className="px-4 py-2 bg-red-900/20 hover:bg-red-900/30 border border-red-800 rounded-lg text-red-400 font-medium transition-colors"
                >
                    Удалить проект
                </button>
            </div>

            {/* Delete Confirmation Modal */}
            <AnimatePresence>
                {showDeleteConfirm && (
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50"
                        onClick={() => !isDeleting && setShowDeleteConfirm(false)}
                    >
                        <motion.div
                            initial={{ scale: 0.95, y: 20 }}
                            animate={{ scale: 1, y: 0 }}
                            exit={{ scale: 0.95, y: 20 }}
                            className="card max-w-md w-full border-red-900/50"
                            onClick={(e) => e.stopPropagation()}
                        >
                            <h3 className="text-xl font-bold text-white mb-2">
                                Удалить проект?
                            </h3>
                            <p className="text-gray-400 mb-2">
                                Вы уверены, что хотите удалить проект <span className="text-white font-semibold">"{project.name}"</span>?
                            </p>
                            <p className="text-red-400 text-sm mb-6">
                                Это действие необратимо. Все коллекции и данные будут удалены.
                            </p>

                            <div className="flex gap-3">
                                <button
                                    type="button"
                                    onClick={() => setShowDeleteConfirm(false)}
                                    className="btn-ghost flex-1"
                                    disabled={isDeleting}
                                >
                                    Отмена
                                </button>
                                <button
                                    onClick={handleDelete}
                                    className="flex-1 px-4 py-2 bg-red-600 hover:bg-red-700 rounded-lg text-white font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                                    disabled={isDeleting}
                                >
                                    {isDeleting ? 'Удаление...' : 'Удалить'}
                                </button>
                            </div>
                        </motion.div>
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    );
}

