import { useState } from 'preact/hooks';
import { motion } from 'framer-motion';

export function Auth({ onLogin, onRegister }) {
    const [isLogin, setIsLogin] = useState(true);
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);

        const result = isLogin 
            ? await onLogin(email, password)
            : await onRegister(email, password);

        setIsLoading(false);

        if (!result.success) {
            setError(result.error);
        }
    };

    return (
        <div className="min-h-screen bg-dark-900 flex items-center justify-center p-4">
            <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className="w-full max-w-md"
            >
                {/* Logo */}
                <motion.div 
                    className="text-center mb-8"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.1 }}
                >
                    <h1 className="text-4xl font-bold gradient-text mb-2">
                        MockingCode
                    </h1>
                    <p className="text-gray-400">
                        Быстрые mock API для разработки
                    </p>
                </motion.div>

                {/* Form Card */}
                <motion.div
                    className="card"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.2 }}
                >
                    <div className="flex gap-2 mb-6">
                        <button
                            onClick={() => {
                                setIsLogin(true);
                                setError('');
                            }}
                            className={`flex-1 py-2 rounded-lg font-medium transition-all ${
                                isLogin 
                                    ? 'bg-primary-600 text-white' 
                                    : 'bg-dark-700 text-gray-400 hover:text-white'
                            }`}
                        >
                            Вход
                        </button>
                        <button
                            onClick={() => {
                                setIsLogin(false);
                                setError('');
                            }}
                            className={`flex-1 py-2 rounded-lg font-medium transition-all ${
                                !isLogin 
                                    ? 'bg-primary-600 text-white' 
                                    : 'bg-dark-700 text-gray-400 hover:text-white'
                            }`}
                        >
                            Регистрация
                        </button>
                    </div>

                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-2">
                                Email
                            </label>
                            <input
                                type="email"
                                value={email}
                                onInput={(e) => setEmail(e.target.value)}
                                className="input"
                                placeholder="your@email.com"
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-2">
                                Пароль
                            </label>
                            <input
                                type="password"
                                value={password}
                                onInput={(e) => setPassword(e.target.value)}
                                className="input"
                                placeholder="••••••••"
                                required
                                minLength={6}
                                disabled={isLoading}
                            />
                        </div>

                        {error && (
                            <motion.div
                                initial={{ opacity: 0, y: -10 }}
                                animate={{ opacity: 1, y: 0 }}
                                className="bg-red-900/20 border border-red-800 rounded-lg p-3 text-red-400 text-sm"
                            >
                                {error}
                            </motion.div>
                        )}

                        <motion.button
                            type="submit"
                            disabled={isLoading}
                            whileHover={{ scale: isLoading ? 1 : 1.02 }}
                            whileTap={{ scale: isLoading ? 1 : 0.98 }}
                            className="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {isLoading ? 'Загрузка...' : (isLogin ? 'Войти' : 'Зарегистрироваться')}
                        </motion.button>
                    </form>
                </motion.div>

                {/* Footer */}
                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.4 }}
                    className="text-center mt-6 text-sm text-gray-500"
                >
                    <p>Бесплатное API для разработчиков</p>
                </motion.div>
            </motion.div>
        </div>
    );
}

