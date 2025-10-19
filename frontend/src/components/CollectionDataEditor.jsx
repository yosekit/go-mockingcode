import { useState, useEffect } from 'preact/hooks';
import { motion, AnimatePresence } from 'framer-motion';
import Editor from '@monaco-editor/react';
import apiClient from '../utils/apiClient';

export function CollectionDataEditor({ apiKey, collection, onClose }) {
    const [documents, setDocuments] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState('');
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [editingDocument, setEditingDocument] = useState(null);
    const [jsonValue, setJsonValue] = useState('{\n  \n}');
    const [isSaving, setIsSaving] = useState(false);
    const [isBulkEdit, setIsBulkEdit] = useState(false);

    useEffect(() => {
        loadDocuments();
    }, [apiKey, collection.name]);

    const loadDocuments = async () => {
        try {
            setIsLoading(true);
            setError('');
            const response = await apiClient.getCollectionData(apiKey, collection.name);
            // Backend теперь возвращает чистый массив
            setDocuments(Array.isArray(response) ? response : []);
        } catch (err) {
            setError(err.message);
        } finally {
            setIsLoading(false);
        }
    };

    const handleCreate = () => {
        setEditingDocument(null);
        setIsBulkEdit(false);
        setJsonValue('{\n  "name": "Example",\n  "value": 123\n}');
        setShowCreateModal(true);
    };

    const handleBulkEdit = () => {
        setEditingDocument(null);
        setIsBulkEdit(true);
        // Формируем JSON массив всех документов (показываем все как есть)
        setJsonValue(JSON.stringify(documents, null, 2));
        setShowCreateModal(true);
    };

    const handleEdit = (doc) => {
        setEditingDocument(doc);
        setIsBulkEdit(false);
        // Убираем id из редактирования (read-only, автогенерируется)
        const { id, ...editableData } = doc;
        setJsonValue(JSON.stringify(editableData, null, 2));
        setShowCreateModal(true);
    };

    const handleSave = async () => {
        try {
            setIsSaving(true);
            setError('');
            
            // Парсим JSON
            const data = JSON.parse(jsonValue);
            
            if (isBulkEdit) {
                // Массовое обновление
                if (!Array.isArray(data)) {
                    setError('Для массового редактирования нужен массив объектов');
                    setIsSaving(false);
                    return;
                }
                
                // Обновляем каждый документ
                const promises = data.map(doc => {
                    const { id, ...updateData } = doc;
                    if (!id) {
                        throw new Error('Каждый документ должен иметь id');
                    }
                    return apiClient.updateDocument(apiKey, collection.name, id, updateData);
                });
                
                await Promise.all(promises);
            } else if (editingDocument) {
                // Обновление одного документа
                await apiClient.updateDocument(apiKey, collection.name, editingDocument.id, data);
            } else {
                // Создание
                await apiClient.createDocument(apiKey, collection.name, data);
            }
            
            setShowCreateModal(false);
            loadDocuments();
        } catch (err) {
            if (err instanceof SyntaxError) {
                setError('Неверный JSON формат');
            } else {
                setError(err.message);
            }
        } finally {
            setIsSaving(false);
        }
    };

    const handleDelete = async (doc) => {
        if (!confirm('Удалить документ?')) return;
        
        try {
            await apiClient.deleteDocument(apiKey, collection.name, doc.id);
            loadDocuments();
        } catch (err) {
            setError(err.message);
        }
    };

    if (isLoading) {
        return (
            <div className="flex items-center justify-center py-12">
                <div className="text-gray-400">Загрузка данных...</div>
            </div>
        );
    }

    return (
        <div className="space-y-4">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h3 className="text-xl font-semibold text-white mb-1">
                        Данные коллекции: {collection.name}
                    </h3>
                    <p className="text-gray-400 text-sm">
                        {documents.length} {documents.length === 1 ? 'документ' : 'документов'}
                    </p>
                </div>
                <div className="flex gap-2">
                    {documents.length > 0 && (
                        <button onClick={handleBulkEdit} className="btn-secondary">
                            Редактировать все
                        </button>
                    )}
                    <button onClick={handleCreate} className="btn-primary">
                        + Добавить документ
                    </button>
                    <button onClick={onClose} className="btn-ghost">
                        Закрыть
                    </button>
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

            {/* Documents List */}
            {documents.length === 0 ? (
                <div className="text-center py-12 border-2 border-dashed border-dark-700 rounded-lg">
                    <p className="text-gray-400 mb-4">Нет данных в коллекции</p>
                    <button onClick={handleCreate} className="btn-primary">
                        Добавить первый документ
                    </button>
                </div>
            ) : (
                <div className="space-y-3">
                    {documents.map((doc, index) => (
                        <motion.div
                            key={doc.id || index}
                            initial={{ opacity: 0, y: 10 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ delay: index * 0.05 }}
                            className="bg-dark-800 border border-dark-700 rounded-lg p-4"
                        >
                            <div className="flex items-start justify-between mb-3">
                                <div className="flex items-center gap-2">
                                    <span className="text-xs text-gray-500 font-mono">
                                        ID: {doc.id}
                                    </span>
                                </div>
                                <div className="flex gap-2">
                                    <button
                                        onClick={() => handleEdit(doc)}
                                        className="text-gray-400 hover:text-primary-400 transition-colors"
                                        title="Редактировать"
                                    >
                                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                        </svg>
                                    </button>
                                    <button
                                        onClick={() => handleDelete(doc)}
                                        className="text-gray-400 hover:text-red-400 transition-colors"
                                        title="Удалить"
                                    >
                                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                        </svg>
                                    </button>
                                </div>
                            </div>
                            
                            {/* JSON Preview */}
                            <pre className="bg-dark-900 rounded p-3 text-xs text-gray-300 overflow-x-auto">
                                {JSON.stringify(doc, null, 2)}
                            </pre>
                        </motion.div>
                    ))}
                </div>
            )}

            {/* Create/Edit Modal */}
            <AnimatePresence>
                {showCreateModal && (
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50"
                        onClick={() => !isSaving && setShowCreateModal(false)}
                    >
                        <motion.div
                            initial={{ scale: 0.95, y: 20 }}
                            animate={{ scale: 1, y: 0 }}
                            exit={{ scale: 0.95, y: 20 }}
                            className="card max-w-3xl w-full"
                            onClick={(e) => e.stopPropagation()}
                        >
                            <h3 className="text-xl font-bold text-white mb-4">
                                {isBulkEdit ? 'Массовое редактирование' : (editingDocument ? 'Редактировать документ' : 'Новый документ')}
                            </h3>

                            <div className="mb-4">
                                <label className="block text-sm font-medium text-gray-300 mb-2">
                                    JSON данные
                                </label>
                                <div className="border border-dark-600 rounded-lg overflow-hidden">
                                    <Editor
                                        height="400px"
                                        defaultLanguage="json"
                                        value={jsonValue}
                                        onChange={(value) => setJsonValue(value || '{}')}
                                        theme="vs-dark"
                                        options={{
                                            minimap: { enabled: false },
                                            fontSize: 14,
                                            lineNumbers: 'on',
                                            scrollBeyondLastLine: false,
                                            automaticLayout: true,
                                            tabSize: 2,
                                        }}
                                    />
                                </div>
                                <p className="text-xs text-gray-500 mt-2">
                                    {isBulkEdit 
                                        ? 'Редактируйте массив документов. Поле id не изменяется, изменяются только данные.'
                                        : 'Введите JSON объект. Поле id генерируется автоматически.'
                                    }
                                </p>
                            </div>

                            {error && (
                                <div className="bg-red-900/20 border border-red-800 rounded-lg p-3 text-red-400 text-sm mb-4">
                                    {error}
                                </div>
                            )}

                            <div className="flex gap-3">
                                <button
                                    type="button"
                                    onClick={() => setShowCreateModal(false)}
                                    className="btn-ghost flex-1"
                                    disabled={isSaving}
                                >
                                    Отмена
                                </button>
                                <button
                                    onClick={handleSave}
                                    className="btn-primary flex-1"
                                    disabled={isSaving}
                                >
                                    {isSaving ? 'Сохранение...' : 'Сохранить'}
                                </button>
                            </div>
                        </motion.div>
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    );
}

