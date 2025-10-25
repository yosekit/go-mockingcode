import { useState } from 'preact/hooks';
import { motion, AnimatePresence } from 'framer-motion';
import apiClient from '../utils/apiClient';

// Доступные типы Faker
const FAKER_TYPES = [
    { value: 'string', label: 'String (текст)' },
    { value: 'number', label: 'Number (число)' },
    { value: 'boolean', label: 'Boolean (true/false)' },
    { value: 'date', label: 'Date (дата)' },
];

const FAKER_FORMATS = {
    string: [
        { value: '', label: 'Случайное слово' },
        { value: 'name', label: 'Имя' },
        { value: 'email', label: 'Email' },
        { value: 'phone', label: 'Телефон' },
        { value: 'username', label: 'Username' },
        { value: 'url', label: 'URL' },
        { value: 'address', label: 'Адрес' },
        { value: 'city', label: 'Город' },
        { value: 'country', label: 'Страна' },
        { value: 'uuid', label: 'UUID' },
    ],
    number: [],
    boolean: [],
    date: [],
};

export function SchemaEditor({ collection, apiKey, onSave, onCancel }) {
    // Инициализируем поля, убеждаясь что id всегда readOnly
    const initFields = () => {
        const existingFields = collection.fields || [];
        
        // Проверяем есть ли уже поле id
        const hasId = existingFields.some(f => f.name === 'id');
        
        if (hasId) {
            // Если есть - помечаем его как readOnly
            return existingFields.map(f => 
                f.name === 'id' ? { ...f, readOnly: true } : f
            );
        } else {
            // Если нет - добавляем
            return [
                { name: 'id', type: 'string', required: true, readOnly: true },
                ...existingFields
            ];
        }
    };

    const [fields, setFields] = useState(initFields());
    const [isSaving, setIsSaving] = useState(false);
    const [error, setError] = useState('');
    const [isEditMode, setIsEditMode] = useState(false);
    const [generateCount, setGenerateCount] = useState(10);

    const addField = () => {
        setFields([...fields, { name: '', type: 'string', format: '', required: false }]);
    };

    const updateField = (index, key, value) => {
        const updated = [...fields];
        updated[index] = { ...updated[index], [key]: value };
        setFields(updated);
    };

    const removeField = (index) => {
        setFields(fields.filter((_, i) => i !== index));
    };

    const handleSaveSchema = async (e) => {
        e.preventDefault();
        
        // Валидация
        const fieldNames = new Set();
        for (const field of fields) {
            if (!field.name.trim()) {
                setError('Все поля должны иметь название');
                return;
            }
            if (fieldNames.has(field.name)) {
                setError(`Дубликат поля: ${field.name}`);
                return;
            }
            fieldNames.add(field.name);
        }

        try {
            setIsSaving(true);
            setError('');
            await onSave(fields, false); // false = не закрывать окно
            setIsEditMode(false); // Возвращаемся к генератору
        } catch (err) {
            // Не показываем ошибку 401 - она обрабатывается глобально
            if (err.message && !err.message.includes('HTTP 401')) {
                setError(err.message);
            }
        } finally {
            setIsSaving(false);
        }
    };

    const handleGenerate = async () => {
        try {
            setIsSaving(true);
            setError('');
            
            // Генерируем данные
            const response = await apiClient.generateDocuments(
                fields, 
                generateCount
            );
            
            // Сохраняем каждый документ через обычный API
            let savedCount = 0;
            for (const doc of response.documents) {
                try {
                    await apiClient.createDocument(apiKey, collection.name, doc);
                    savedCount++;
                } catch (err) {
                    console.error('Failed to save document:', err);
                }
            }
            
            // Показываем результат
            alert(`Сгенерировано ${response.count} документов, сохранено ${savedCount}!`);
            
            // Вызываем callback для обновления данных в родительском компоненте
            if (onSave) {
                onSave();
            }
            
            // Закрываем окно
            onCancel();
        } catch (err) {
            // Не показываем ошибку 401 - она обрабатывается глобально
            if (err.message && !err.message.includes('HTTP 401')) {
                setError(err.message);
            }
        } finally {
            setIsSaving(false);
        }
    };

    return (
        <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50"
            onClick={onCancel}
        >
            <motion.div
                initial={{ scale: 0.95, y: 20 }}
                animate={{ scale: 1, y: 0 }}
                exit={{ scale: 0.95, y: 20 }}
                className="card max-w-3xl w-full max-h-[90vh] overflow-y-auto"
                onClick={(e) => e.stopPropagation()}
            >
                <div className="flex items-center justify-between mb-4">
                    <h3 className="text-xl font-bold text-white">
                        {isEditMode ? 'Редактор схемы' : 'Генератор данных'}: {collection.name}
                    </h3>
                    <div className="flex gap-2">
                        {!isEditMode && (
                            <button
                                type="button"
                                onClick={() => setIsEditMode(true)}
                                className="btn-ghost text-sm"
                            >
                                Редактировать схему
                            </button>
                        )}
                        <button
                            type="button"
                            onClick={onCancel}
                            className="p-2 text-gray-400 hover:text-white transition-colors"
                            title="Закрыть"
                        >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                </div>

                <form onSubmit={isEditMode ? handleSaveSchema : (e) => { e.preventDefault(); handleGenerate(); }}>
                    <div className="space-y-3 mb-6">
                        {fields.map((field, index) => (
                            <div key={index} className={`rounded-lg p-3 ${field.readOnly ? 'bg-dark-700/50 opacity-60' : 'bg-dark-700/50'}`}>
                                {field.readOnly && (
                                    <div className="text-xs text-gray-500 mb-2">
                                        🔒 Автоматическое поле (read-only)
                                    </div>
                                )}
                                <div className="grid grid-cols-12 gap-3 items-start">
                                    {/* Название поля */}
                                    <div className="col-span-4">
                                        <input
                                            type="text"
                                            value={field.name}
                                            onInput={(e) => updateField(index, 'name', e.target.value)}
                                            className="input text-sm"
                                            placeholder="field_name"
                                            disabled={!isEditMode || field.readOnly || isSaving}
                                            required
                                        />
                                    </div>

                                    {/* Тип */}
                                    <div className="col-span-3">
                                        <select
                                            value={field.type}
                                            onChange={(e) => updateField(index, 'type', e.target.value)}
                                            className="input text-sm"
                                            disabled={!isEditMode || field.readOnly || isSaving}
                                        >
                                            {FAKER_TYPES.map(t => (
                                                <option key={t.value} value={t.value}>
                                                    {t.label}
                                                </option>
                                            ))}
                                        </select>
                                    </div>

                                    {/* Формат (для string) */}
                                    <div className="col-span-4">
                                        {field.type === 'string' && (
                                                <select
                                                    value={field.format || ''}
                                                    onChange={(e) => updateField(index, 'format', e.target.value)}
                                                    className="input text-sm"
                                                    disabled={!isEditMode || field.readOnly || isSaving}
                                            >
                                                {FAKER_FORMATS.string.map(f => (
                                                    <option key={f.value} value={f.value}>
                                                        {f.label}
                                                    </option>
                                                ))}
                                            </select>
                                        )}
                                        {field.type === 'number' && (
                                            <input
                                                type="text"
                                                value={field.format || ''}
                                                onInput={(e) => updateField(index, 'format', e.target.value)}
                                                className="input text-sm"
                                                placeholder="min:0,max:100"
                                                disabled={!isEditMode || field.readOnly || isSaving}
                                            />
                                        )}
                                    </div>

                                    {/* Удалить */}
                                    <div className="col-span-1 flex justify-end">
                                        {isEditMode && !field.readOnly && (
                                            <button
                                                type="button"
                                                onClick={() => removeField(index)}
                                                className="text-gray-400 hover:text-red-400 transition-colors"
                                                disabled={isSaving}
                                            >
                                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                                </svg>
                                            </button>
                                        )}
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>

                    {isEditMode && (
                        <button
                            type="button"
                            onClick={addField}
                            className="btn-ghost w-full mb-6"
                            disabled={isSaving}
                        >
                            + Добавить поле
                        </button>
                    )}

                    {!isEditMode && (
                        <div className="mb-6">
                            <label className="block text-sm font-medium text-gray-300 mb-2">
                                Количество документов для генерации
                            </label>
                            <input
                                type="number"
                                value={generateCount}
                                onInput={(e) => setGenerateCount(parseInt(e.target.value) || 1)}
                                className="input"
                                min="1"
                                max="100"
                                disabled={isSaving}
                            />
                        </div>
                    )}

                    {error && (
                        <div className="bg-red-900/20 border border-red-800 rounded-lg p-3 text-red-400 text-sm mb-4">
                            {error}
                        </div>
                    )}

                    <div className="flex gap-3">
                        {isEditMode && (
                            <button
                                type="button"
                                onClick={() => {
                                    setIsEditMode(false);
                                    setFields(initFields()); // Сбрасываем изменения
                                    setError('');
                                }}
                                className="btn-ghost flex-1"
                                disabled={isSaving}
                            >
                                Отмена
                            </button>
                        )}
                        {!isEditMode && (
                            <button
                                type="button"
                                onClick={onCancel}
                                className="btn-ghost flex-1"
                                disabled={isSaving}
                            >
                                Закрыть
                            </button>
                        )}
                        <button
                            type="submit"
                            className="btn-primary flex-1"
                            disabled={isSaving}
                        >
                            {isSaving 
                                ? (isEditMode ? 'Сохранение...' : 'Генерация...')
                                : (isEditMode ? 'Сохранить схему' : 'Сгенерировать')
                            }
                        </button>
                    </div>
                </form>
            </motion.div>
        </motion.div>
    );
}

