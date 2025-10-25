import { useState, useEffect } from 'preact/hooks';
import { motion, AnimatePresence } from 'framer-motion';
import apiClient from '../utils/apiClient';
import { CollectionDataEditor } from './CollectionDataEditor';

export function Collections({ projectId, apiKey, limits }) {
    const [collections, setCollections] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState('');
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [newCollectionName, setNewCollectionName] = useState('');
    const [newCollectionDescription, setNewCollectionDescription] = useState('');
    const [isCreating, setIsCreating] = useState(false);
    const [selectedCollection, setSelectedCollection] = useState(null);
    const [documentCounts, setDocumentCounts] = useState({});

    // Используем лимиты из пропсов или дефолтные значения
    const effectiveLimits = limits || {
        max_collections_per_project: 20,
        max_documents_per_collection: 500
    };

    useEffect(() => {
        loadCollections();
    }, [projectId]);

    const loadCollections = async () => {
        try {
            setIsLoading(true);
            const data = await apiClient.getCollections(projectId);
            setCollections(data || []);
            
            // Загружаем количество документов для каждой коллекции
            if (data && data.length > 0) {
                loadDocumentCounts(data);
            }
        } catch (err) {
            // Не показываем ошибку 401 - она обрабатывается глобально
            if (err.message && !err.message.includes('HTTP 401')) {
                setError(err.message);
            }
        } finally {
            setIsLoading(false);
        }
    };

    const loadDocumentCounts = async (collections) => {
        const counts = {};
        
        for (const collection of collections) {
            try {
                const documents = await apiClient.getCollectionData(apiKey, collection.name);
                counts[collection.id] = Array.isArray(documents) ? documents.length : 0;
            } catch (err) {
                // Если коллекция пустая или ошибка - считаем 0
                counts[collection.id] = 0;
            }
        }
        
        setDocumentCounts(counts);
    };

    const handleCreateCollection = async (e) => {
        e.preventDefault();
        if (!newCollectionName.trim()) return;

        // Проверяем лимит коллекций
        if (collections.length >= effectiveLimits.max_collections_per_project) {
            setError(`Достигнут лимит коллекций: ${effectiveLimits.max_collections_per_project}`);
            return;
        }

        try {
            setIsCreating(true);
            const collection = await apiClient.createCollection(projectId, {
                name: newCollectionName,
                description: newCollectionDescription,
                // Backend требует хотя бы одно поле
                fields: [
                    {
                        name: 'id',
                        type: 'string',
                        required: true,
                    }
                ],
            });
            setCollections([...collections, collection]);
            setDocumentCounts({...documentCounts, [collection.id]: 0});
            setNewCollectionName('');
            setNewCollectionDescription('');
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

    const handleDeleteCollection = async (collectionId) => {
        if (!confirm('Удалить коллекцию? Все данные будут потеряны.')) return;

        try {
            await apiClient.deleteCollection(projectId, collectionId);
            setCollections(collections.filter(c => c.id !== collectionId));
            
            // Удаляем счетчик документов для удаленной коллекции
            const newCounts = { ...documentCounts };
            delete newCounts[collectionId];
            setDocumentCounts(newCounts);
        } catch (err) {
            setError(err.message);
        }
    };


    if (isLoading) {
        return (
            <div className="flex items-center justify-center py-8">
                <div className="text-gray-400">Загрузка коллекций...</div>
            </div>
        );
    }

    // Если выбрана коллекция - показываем редактор данных
    if (selectedCollection) {
        return (
            <CollectionDataEditor
                apiKey={apiKey}
                collection={selectedCollection}
                projectId={projectId}
                limits={effectiveLimits}
                onClose={() => {
                    setSelectedCollection(null);
                    // Обновляем счетчики документов при возврате
                    loadDocumentCounts(collections);
                }}
                onDataUpdate={() => {
                    // Обновляем счетчики документов после изменения данных
                    loadDocumentCounts(collections);
                }}
            />
        );
    }

    return (
        <div>
            {/* Header */}
            <div className="flex justify-between items-center mb-4">
                <div>
                    <h3 className="text-xl font-semibold text-white">Коллекции</h3>
                    <p className="text-gray-400 text-sm mt-1">
                        Создавайте коллекции для хранения mock данных
                        <span className="text-gray-500 ml-2">
                            ({collections.length}/{effectiveLimits.max_collections_per_project})
                        </span>
                    </p>
                </div>
                <motion.button
                    onClick={() => setShowCreateModal(true)}
                    whileHover={{ scale: collections.length < effectiveLimits.max_collections_per_project ? 1.05 : 1 }}
                    whileTap={{ scale: collections.length < effectiveLimits.max_collections_per_project ? 0.95 : 1 }}
                    className={`btn-primary ${collections.length >= effectiveLimits.max_collections_per_project ? 'opacity-50 cursor-not-allowed' : ''}`}
                    disabled={collections.length >= effectiveLimits.max_collections_per_project}
                >
                    + Новая коллекция
                </motion.button>
            </div>

            {/* Error */}
            {error && (
                <motion.div
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="bg-red-900/20 border border-red-800 rounded-lg p-3 text-red-400 text-sm mb-4"
                >
                    {error}
                </motion.div>
            )}

            {/* Collections List */}
            {collections.length === 0 ? (
                <div className="text-center py-12 border-2 border-dashed border-dark-700 rounded-lg">
                    <div className="text-gray-500 mb-4">
                        <svg className="w-12 h-12 mx-auto mb-3 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                        </svg>
                    </div>
                    <h4 className="text-lg font-semibold text-white mb-2">
                        Нет коллекций
                    </h4>
                    <p className="text-gray-400 mb-4">
                        Создайте первую коллекцию для хранения данных
                    </p>
                    <button
                        onClick={() => setShowCreateModal(true)}
                        className="btn-primary"
                    >
                        Создать коллекцию
                    </button>
                </div>
            ) : (
                <div className="space-y-3">
                    {collections.map((collection, index) => (
                        <motion.div
                            key={collection.id}
                            initial={{ opacity: 0, y: 10 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ delay: index * 0.05 }}
                            className="bg-dark-800 border border-dark-700 rounded-lg p-4 hover:border-primary-600/50 transition-colors"
                        >
                            <div className="flex items-start justify-between">
                                <div className="flex-1">
                                    <div className="flex items-center gap-3 mb-2">
                                        <h4 className="text-lg font-semibold text-white">
                                            /{collection.name}
                                        </h4>
                                        <button
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                const apiUrl = `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/${apiKey}/${collection.name}`;
                                                window.open(apiUrl, '_blank');
                                            }}
                                            className="p-1.5 text-gray-400 hover:text-blue-400 transition-all duration-200 hover:scale-110"
                                            title="Открыть API коллекции"
                                        >
                                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                                            </svg>
                                        </button>
                                    </div>
                                    
                                    {collection.description && (
                                        <p className="text-gray-400 text-sm mb-3">
                                            {collection.description}
                                        </p>
                                    )}

                                    <div className="flex items-center gap-4 text-xs text-gray-500">
                                        <div className="flex items-center gap-4">
                                            <div className="flex items-center gap-3">
                                                <div className="w-16 h-1 rounded-md bg-gray-700">
                                                <div 
                                                    style={{ width: `${Math.min((documentCounts[collection.id] || 0) / effectiveLimits.max_documents_per_collection * 100, 100)}%` }}
                                                    className="h-full bg-lime-500 rounded-md transition-all duration-300"
                                                ></div>
                                            </div>
                                            <span className="text-gray-500">
                                                {documentCounts[collection.id] || 0} / {effectiveLimits.max_documents_per_collection || 'undefined'}
                                            </span>
                                            </div>
                                            
                                            {collection.schema && (
                                                <div className="flex items-center gap-1">
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                                    </svg>
                                                    <span>Схема определена</span>
                                                </div>
                                            )}
                                        </div>
                                        
                                        <span className={`px-2 py-0.5 rounded text-xs font-medium ${
                                            collection.is_active 
                                                ? 'bg-green-900/30 text-green-400 border border-green-800' 
                                                : 'bg-gray-900/30 text-gray-400 border border-gray-800'
                                        }`}>
                                            {collection.is_active ? 'Активна' : 'Неактивна'}
                                        </span>
                                    </div>
                                </div>

                                <div className="flex items-center gap-2">
                                    <button
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            setSelectedCollection(collection);
                                        }}
                                        className="px-3 py-1.5 bg-primary-600 hover:bg-primary-700 rounded text-white text-sm font-medium transition-colors"
                                    >
                                        Данные
                                    </button>
                                    <button
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            handleDeleteCollection(collection.id);
                                        }}
                                        className="p-2 text-gray-400 hover:text-red-400 transition-colors"
                                        title="Удалить коллекцию"
                                    >
                                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                        </svg>
                                    </button>
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
                        onClick={() => !isCreating && setShowCreateModal(false)}
                    >
                        <motion.div
                            initial={{ scale: 0.95, y: 20 }}
                            animate={{ scale: 1, y: 0 }}
                            exit={{ scale: 0.95, y: 20 }}
                            className="card max-w-md w-full"
                            onClick={(e) => e.stopPropagation()}
                        >
                            <h3 className="text-xl font-bold text-white mb-4">
                                Новая коллекция
                            </h3>
                            
                            <form onSubmit={handleCreateCollection}>
                                <div className="mb-4">
                                    <label className="block text-sm font-medium text-gray-300 mb-2">
                                        Название коллекции
                                    </label>
                                    <input
                                        type="text"
                                        value={newCollectionName}
                                        onInput={(e) => setNewCollectionName(e.target.value)}
                                        className="input"
                                        placeholder="users"
                                        required
                                        disabled={isCreating}
                                        autoFocus
                                    />
                                    <p className="text-xs text-gray-500 mt-1">
                                        Используйте латиницу, без пробелов
                                    </p>
                                </div>

                                <div className="mb-4">
                                    <label className="block text-sm font-medium text-gray-300 mb-2">
                                        Описание (опционально)
                                    </label>
                                    <textarea
                                        value={newCollectionDescription}
                                        onInput={(e) => setNewCollectionDescription(e.target.value)}
                                        className="input resize-none"
                                        placeholder="Список пользователей"
                                        rows={2}
                                        disabled={isCreating}
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

