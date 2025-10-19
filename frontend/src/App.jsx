import { useState } from 'preact/hooks'
import { motion } from 'framer-motion'

export function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="min-h-screen bg-dark-900 flex items-center justify-center">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center"
      >
        <motion.h1 
          className="text-4xl font-bold gradient-text mb-8"
          whileHover={{ scale: 1.05 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          MockingCode
        </motion.h1>
        
        <motion.p 
          className="text-gray-400 text-lg mb-8"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.2 }}
        >
          API Mock Generator
        </motion.p>

        <motion.div 
          className="card max-w-md mx-auto"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <h2 className="text-xl font-semibold mb-4">Test Counter</h2>
          <p className="text-gray-300 mb-4">
            Current count: <span className="text-primary-400 font-mono">{count}</span>
          </p>
          
          <motion.button
            onClick={() => setCount(count + 1)}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            className="btn-primary"
          >
            Increment
          </motion.button>
        </motion.div>

        <motion.div 
          className="mt-8 text-sm text-gray-500"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
        >
          <p>✅ Preact + Vite + Tailwind CSS + Framer Motion</p>
        </motion.div>
      </motion.div>
    </div>
  )
}