import React, { useState, useEffect } from 'react';

// Components
import Header from './components/Header';
import WorkflowList from './components/WorkflowList';
import RunList from './components/RunList';
import RunDetails from './components/RunDetails';

// Services
import api from './services/api';

export default function App() {
  const [workflows, setWorkflows] = useState([]);
  const [runs, setRuns] = useState([]);
  const [selectedRun, setSelectedRun] = useState(null);
  const [loading, setLoading] = useState(false);
  const [uploadStatus, setUploadStatus] = useState('');

  // Fetch data on mount and setup polling
  useEffect(() => {
    fetchWorkflows();
    fetchRuns();
    const interval = setInterval(fetchRuns, 3000);
    return () => clearInterval(interval);
  }, []);

  const fetchWorkflows = async () => {
    try {
      const data = await api.getWorkflows();
      setWorkflows(data || []);
    } catch (err) {
      console.error('Failed to fetch workflows:', err);
    }
  };

  const fetchRuns = async () => {
    try {
      const data = await api.getRuns();
      setRuns(data);
    } catch (err) {
      console.error('Failed to fetch runs:', err);
    }
  };

  const uploadWorkflow = async (e) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setLoading(true);
    setUploadStatus('');

    try {
      await api.uploadWorkflow(file);
      setUploadStatus('✓ Workflow uploaded successfully');
      fetchWorkflows();
    } catch (err) {
      setUploadStatus('✗ ' + err.message);
    } finally {
      setLoading(false);
      // Clear file input
      e.target.value = '';
    }
  };

  const triggerWorkflow = async (name) => {
    setLoading(true);
    try {
      await api.triggerWorkflow(name);
      fetchRuns();
    } catch (err) {
      console.error('Failed to trigger workflow:', err);
      alert('Failed to trigger workflow: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const viewRunDetails = async (runId) => {
    try {
      const data = await api.getRun(runId);
      setSelectedRun(data);
    } catch (err) {
      console.error('Failed to fetch run details:', err);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto p-6">
        <Header />

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
          <div className="lg:col-span-1">
            <WorkflowList
              workflows={workflows}
              onUpload={uploadWorkflow}
              onTrigger={triggerWorkflow}
              uploadStatus={uploadStatus}
              loading={loading}
            />
          </div>

          <div className="lg:col-span-2">
            <RunList
              runs={runs}
              onSelectRun={viewRunDetails}
              onRefresh={fetchRuns}
            />
          </div>
        </div>

        {selectedRun && <RunDetails run={selectedRun} />}
      </div>
    </div>
  );
}