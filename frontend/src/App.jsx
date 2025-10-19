import { useState } from 'preact/hooks'
import { useAuth } from './hooks/useAuth'
import { Auth } from './components/Auth'
import { Projects } from './components/Projects'
import { ProjectDetail } from './components/ProjectDetail'

export function App() {
  const { isAuthenticated, isLoading, login, register, logout } = useAuth();
  const [selectedProject, setSelectedProject] = useState(null);

  const handleProjectUpdated = (updatedProject) => {
    setSelectedProject(updatedProject);
  };

  const handleProjectDeleted = () => {
    setSelectedProject(null);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-dark-900 flex items-center justify-center">
        <div className="text-gray-400">Загрузка...</div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Auth onLogin={login} onRegister={register} />;
  }

  return (
    <div className="min-h-screen bg-dark-900 p-8">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div className="flex items-center gap-4">
            <h1 className="text-3xl font-bold gradient-text">MockingCode</h1>
            {selectedProject && (
              <button
                onClick={() => setSelectedProject(null)}
                className="text-gray-400 hover:text-white transition-colors text-sm"
              >
                ← Назад к проектам
              </button>
            )}
          </div>
          <button onClick={logout} className="btn-secondary">
            Выйти
          </button>
        </div>
        
        {/* Content */}
        {selectedProject ? (
          <ProjectDetail 
            project={selectedProject}
            onBack={() => setSelectedProject(null)}
            onProjectUpdated={handleProjectUpdated}
            onProjectDeleted={handleProjectDeleted}
          />
        ) : (
          <Projects onSelectProject={setSelectedProject} />
        )}
      </div>
    </div>
  );
}