import { useState, useEffect } from 'preact/hooks';
import { motion, AnimatePresence } from 'framer-motion';
import Editor from '@monaco-editor/react';
import apiClient from '../utils/apiClient';
import { SchemaEditor } from './SchemaEditor';
import { generateDefaultDocument, formatDocumentAsJson } from '../utils/schemaDefaults';

export function CollectionDataEditor({ apiKey, collection, onClose, projectId, onDataUpdate, limits }) {
    const [documents, setDocuments] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState('');
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [editingDocument, setEditingDocument] = useState(null);
    const [jsonValue, setJsonValue] = useState('{\n  \n}');
    const [isSaving, setIsSaving] = useState(false);
    const [isBulkEdit, setIsBulkEdit] = useState(false);
    const [showSchemaEditor, setShowSchemaEditor] = useState(false);
    const [showFlushConfirm, setShowFlushConfirm] = useState(false);
    const [currentCollection, setCurrentCollection] = useState(collection);

    // Используем лимиты из пропсов или дефолтные значения
    const effectiveLimits = limits || {
        max_collections_per_project: 20,
        max_documents_per_collection: 500
    };

    useEffect(() => {
        loadDocuments();
    }, [apiKey, currentCollection.name]);

    const loadDocuments = async () => {
        try {
            setIsLoading(true);
            setError('');
            const response = await apiClient.getCollectionData(apiKey, currentCollection.name);
            // Backend теперь возвращает чистый массив
            const docs = Array.isArray(response) ? response : [];
            
            // Сортируем документы по ID по возрастанию
            const sortedDocs = docs.sort((a, b) => {
                const idA = parseInt(a.id) || 0;
                const idB = parseInt(b.id) || 0;
                return idA - idB;
            });
            
            setDocuments(sortedDocs);
        } catch (err) {
            // Не показываем ошибку 401 - она обрабатывается глобально
            if (err.message && !err.message.includes('HTTP 401')) {
                setError(err.message);
            }
        } finally {
            setIsLoading(false);
        }
    };

    const handleCreate = () => {
        setEditingDocument(null);
        setIsBulkEdit(false);
        
        // Генерируем дефолтный JSON на основе схемы полей коллекции
        const fields = currentCollection.fields || [];
        const defaultDocument = generateDefaultDocument(fields);
        const jsonString = formatDocumentAsJson(defaultDocument);
        
        setJsonValue(jsonString);
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
        
        // Если у документа нет полей из схемы, добавляем дефолтные значения
        const fields = currentCollection.fields || [];
        const defaultDocument = generateDefaultDocument(fields);
        
        // Объединяем существующие данные с дефолтными значениями для недостающих полей
        const mergedData = { ...defaultDocument, ...editableData };
        
        setJsonValue(formatDocumentAsJson(mergedData));
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
                    return apiClient.updateDocument(apiKey, currentCollection.name, id, updateData);
                });
                
                await Promise.all(promises);
            } else if (editingDocument) {
                // Обновление одного документа
                await apiClient.updateDocument(apiKey, currentCollection.name, editingDocument.id, data);
            } else {
                // Создание
                await apiClient.createDocument(apiKey, currentCollection.name, data);
            }
            
            setShowCreateModal(false);
            loadDocuments();
            
            // Уведомляем родительский компонент об обновлении данных
            if (onDataUpdate) {
                onDataUpdate();
            }
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

    const handleFlushCollection = async () => {
        try {
            await apiClient.flushCollection(apiKey, currentCollection.name);
            loadDocuments();
            setShowFlushConfirm(false);
            
            // Уведомляем родительский компонент об обновлении данных
            if (onDataUpdate) {
                onDataUpdate();
            }
        } catch (err) {
            setError(err.message);
        }
    };

    const handleSaveSchema = async (fields, shouldClose = false) => {
        try {
            const updated = await apiClient.updateCollection(projectId, currentCollection.id, {
                fields: fields,
            });
            setCurrentCollection({ ...currentCollection, fields: updated.fields });
            
            // Обновляем документы после изменения схемы
            loadDocuments();
            
            // Уведомляем родительский компонент об обновлении данных
            if (onDataUpdate) {
                onDataUpdate();
            }
            
            if (shouldClose) {
                setShowSchemaEditor(false);
            }
        } catch (err) {
            throw err;
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
                <div className="flex items-center gap-3">
                    <h3 className="text-xl font-semibold text-white">
                        {currentCollection.name}
                    </h3>
                    <button
                        onClick={onClose}
                        className="text-gray-400 hover:text-white transition-colors text-sm"
                    >
                        ← Выйти из редактирования
                    </button>
                </div>
                <div className="flex items-center gap-3">
                    <div className="flex items-center gap-3">
                        <div className="w-16 h-1 rounded-md bg-gray-700">
                            <div 
                                style={{ width: `${Math.min(documents.length / effectiveLimits.max_documents_per_collection * 100, 100)}%` }}
                                className="h-full bg-lime-500 rounded-md transition-all duration-300"
                            ></div>
                        </div>
                        <p className="text-sm text-gray-500">
                            {documents.length} / {effectiveLimits.max_documents_per_collection || 'undefined'}
                        </p>
                    </div>
                    <button
                        onClick={() => setShowSchemaEditor(true)}
                        className="px-4 py-2 border border-white rounded-lg text-white font-medium hover:bg-white hover:text-dark-900 transition-all"
                    >
                        Генератор
                    </button>
                    {documents.length > 0 && (
                        <button
                            onClick={handleBulkEdit}
                            className="p-2 text-gray-400 hover:text-white transition-colors"
                            title="Редактировать все"
                        >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                            </svg>
                        </button>
                    )}
                    <button
                        onClick={handleCreate}
                        className="p-2 text-gray-400 hover:text-white transition-colors"
                        title="Добавить документ"
                    >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                        </svg>
                    </button>
                    <button
                        onClick={() => setShowFlushConfirm(true)}
                        className="p-2 text-gray-400 hover:text-red-400 transition-colors"
                        title="Очистить коллекцию"
                    >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
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
                                        onClick={async () => {
                                            try {
                                                await apiClient.deleteDocument(apiKey, currentCollection.name, doc.id);
                                                loadDocuments();
                                                // Уведомляем родительский компонент об обновлении данных
                                                if (onDataUpdate) {
                                                    onDataUpdate();
                                                }
                                            } catch (err) {
                                                setError(err.message);
                                            }
                                        }}
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
                                        : 'JSON объект сгенерирован на основе схемы полей коллекции. Поле id генерируется автоматически.'
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

            {/* Schema Editor Modal */}
            <AnimatePresence>
                {showSchemaEditor && (
                    <SchemaEditor
                        collection={currentCollection}
                        apiKey={apiKey}
                        onSave={handleSaveSchema}
                        onCancel={() => setShowSchemaEditor(false)}
                    />
                )}
            </AnimatePresence>

            {/* Flush Collection Confirmation Modal */}
            <AnimatePresence>
                {showFlushConfirm && (
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50"
                        onClick={() => setShowFlushConfirm(false)}
                    >
                        <motion.div
                            initial={{ scale: 0.95, y: 20 }}
                            animate={{ scale: 1, y: 0 }}
                            exit={{ scale: 0.95, y: 20 }}
                            className="card max-w-md w-full"
                            onClick={(e) => e.stopPropagation()}
                        >
                            <h3 className="text-xl font-bold text-white mb-4">
                                Очистить коллекцию
                            </h3>
                            
                            <p className="text-gray-300 mb-6">
                                Вы уверены, что хотите удалить все документы из коллекции "{currentCollection.name}"? 
                                Это действие нельзя отменить.
                            </p>

                            <div className="flex gap-3">
                                <button
                                    type="button"
                                    onClick={() => setShowFlushConfirm(false)}
                                    className="btn-ghost flex-1"
                                >
                                    Отмена
                                </button>
                                <button
                                    onClick={handleFlushCollection}
                                    className="btn-primary flex-1 bg-red-600 hover:bg-red-700"
                                >
                                    Очистить
                                </button>
                            </div>
                        </motion.div>
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    );
}

