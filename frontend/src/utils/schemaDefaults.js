// Утилита для генерации дефолтных значений на основе схемы полей

// Дефолтные значения для каждого типа поля
const DEFAULT_VALUES = {
    string: {
        '': 'Пример текста',
        'name': 'Иван Иванов',
        'email': 'ivan@example.com',
        'phone': '+7 (999) 123-45-67',
        'username': 'ivan_user',
        'url': 'https://example.com',
        'address': 'ул. Примерная, д. 1',
        'city': 'Москва',
        'country': 'Россия',
        'uuid': '550e8400-e29b-41d4-a716-446655440000'
    },
    number: {
        '': 42,
        'min:0,max:100': 50,
        'min:1,max:10': 5
    },
    boolean: {
        '': true
    },
    date: {
        '': '2024-01-15T10:30:00Z'
    }
};

/**
 * Генерирует дефолтное значение для поля на основе его типа и формата
 * @param {Object} field - Объект поля с полями name, type, format
 * @returns {*} Дефолтное значение
 */
export function getDefaultValueForField(field) {
    const { type, format = '' } = field;

    // Для поля id всегда возвращаем null (будет автогенерироваться)
    if (field.name === 'id') {
        return null;
    }

    const typeDefaults = DEFAULT_VALUES[type];
    if (!typeDefaults) {
        return null;
    }

    // Если есть конкретный формат, используем его
    if (format && typeDefaults[format] !== undefined) {
        return typeDefaults[format];
    }

    // Иначе используем дефолтное значение для типа
    return typeDefaults[''] || null;
}

/**
 * Генерирует JSON объект с дефолтными значениями на основе схемы полей
 * @param {Array} fields - Массив полей схемы
 * @returns {Object} JSON объект с дефолтными значениями
 */
export function generateDefaultDocument(fields) {
    const document = {};

    fields.forEach(field => {
        if (field.name && field.name.trim()) {
            const defaultValue = getDefaultValueForField(field);
            if (defaultValue !== null) {
                document[field.name] = defaultValue;
            }
        }
    });

    return document;
}

/**
 * Форматирует JSON объект в строку для отображения в редакторе
 * @param {Object} document - JSON объект
 * @returns {string} Отформатированная JSON строка
 */
export function formatDocumentAsJson(document) {
    return JSON.stringify(document, null, 2);
}
