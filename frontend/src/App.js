import React, { useState, useEffect } from 'react';
import { Play, Upload, RefreshCw, CheckCircle, XCircle, Clock, Loader } from 'lucide-react';

const API_URL = 'http://localhost:8080/api';

function App() {
  const [workflows, setWorkflows] = useState([]);
  const [runs, setRuns] = useState([]);
  const [selectedRun, setSelectedRun] = useState(null);
  const [loading, setLoading] = useState(false);
  const [uploadStatus, setUploadStatus] = useState('');

  useEffect(() => {
    fetchWorkflows();
    fetchRuns();
    const interval = setInterval(fetchRuns, 3000);
    return () => clearInterval(interval);
  }, []);

  const fetchWorkflows = async () => {
    try {
      const res = await fetch(`${API_URL}/workflows`);
      const data = await res.json();
      setWorkflows(data || []);
    } catch (err) {
      console.error('Failed to fetch workflows:', err);
    }
  };

  const fetchRuns = async () => {
    try {
      const res = await fetch(`${API_URL}/runs`);
      const data = await res.json();
      setRuns((data || []).sort((a, b) => new Date(b.started_at) - new Date(a.started_at)));
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
      const text = await file.text();
      const res = await fetch(`${API_URL}/workflows`, {
        method: 'POST',
        body: text,
        headers: { 'Content-Type': 'text/yaml' }
      });
      
      if (res.ok) {
        setUploadStatus('✓ Workflow uploaded successfully');
        fetchWorkflows();
      } else {
        setUploadStatus('✗ Failed to upload workflow');
      }
    } catch (err) {
      setUploadStatus('✗ Error: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const triggerWorkflow = async (name) => {
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/workflows/${name}/trigger`, {
        method: 'POST'
      });
      if (res.ok) {
        fetchRuns();
      }
    } catch (err) {
      console.error('Failed to trigger workflow:', err);
    } finally {
      setLoading(false);
    }
  };

  const viewRunDetails = async (runId) => {
    try {
      const res = await fetch(`${API_URL}/runs/${runId}`);
      const data = await res.json();
      setSelectedRun(data);
    } catch (err) {
      console.error('Failed to fetch run details:', err);
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'success':
        return <CheckCircle className="w-5 h-5 text-green-500" />;
      case 'failed':
        return <XCircle className="w-5 h-5 text-red-500" />;
      case 'running':
        return <Loader className="w-5 h-5 text-blue-500 animate-spin" />;
      default:
        return <Clock className="w-5 h-5 text-gray-400" />;
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'success':
        return 'bg-green-100 text-green-800';
      case 'failed':
        return 'bg-red-100 text-red-800';
      case 'running':
        return 'bg-blue-100 text-blue-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto p-6">
        <header className="mb-8">
          <div className="flex items-center gap-4 mb-2">
            <div className="w-12 h-12 bg-gradient-to-br from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-xl">G</span>
            </div>
            <div>
              <h1 className="text-4xl font-bold text-gray-900">Gantry</h1>
              <p className="text-gray-600">Lightweight CI/CD Platform</p>
            </div>
          </div>
        </header>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
          <div className="lg:col-span-1 bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold mb-4">Workflows</h2>
            
            <div className="mb-4">
              <label className="flex items-center justify-center px-4 py-3 bg-blue-500 text-white rounded-lg cursor-pointer hover:bg-blue-600 transition">
                <Upload className="w-5 h-5 mr-2" />
                Upload Workflow YAML
                <input
                  type="file"
                  accept=".yml,.yaml"
                  onChange={uploadWorkflow}
                  className="hidden"
                />
              </label>
              {uploadStatus && (
                <p className={`mt-2 text-sm ${uploadStatus.startsWith('✓') ? 'text-green-600' : 'text-red-600'}`}>
                  {uploadStatus}
                </p>
              )}
            </div>

            <div className="space-y-2">
              {workflows.length === 0 ? (
                <p className="text-gray-500 text-sm">No workflows yet. Upload one to get started!</p>
              ) : (
                workflows.map((wf) => (
                  <div key={wf.name} className="flex items-center justify-between p-3 bg-gray-50 rounded border">
                    <span className="font-medium text-sm">{wf.name}</span>
                    <button
                      onClick={() => triggerWorkflow(wf.name)}
                      disabled={loading}
                      className="p-2 text-blue-600 hover:bg-blue-50 rounded disabled:opacity-50"
                    >
                      <Play className="w-4 h-4" />
                    </button>
                  </div>
                ))
              )}
            </div>
          </div>

          <div className="lg:col-span-2 bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold">Recent Runs</h2>
              <button
                onClick={fetchRuns}
                className="p-2 text-gray-600 hover:bg-gray-100 rounded"
              >
                <RefreshCw className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-3">
              {runs.length === 0 ? (
                <p className="text-gray-500 text-sm">No runs yet. Trigger a workflow to see runs here.</p>
              ) : (
                runs.map((run) => (
                  <div
                    key={run.id}
                    onClick={() => viewRunDetails(run.id)}
                    className="p-4 border rounded-lg hover:bg-gray-50 cursor-pointer transition"
                  >
                    <div className="flex items-center justify-between mb-2">
                      <div className="flex items-center gap-3">
                        {getStatusIcon(run.status)}
                        <div>
                          <h3 className="font-semibold">{run.workflow_name}</h3>
                          <p className="text-sm text-gray-500">{run.id}</p>
                        </div>
                      </div>
                      <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(run.status)}`}>
                        {run.status}
                      </span>
                    </div>
                    <div className="text-xs text-gray-500">
                      Started: {new Date(run.started_at).toLocaleString()}
                      {run.completed_at && ` • Completed: ${new Date(run.completed_at).toLocaleString()}`}
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>

        {selectedRun && (
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold mb-4">Run Details: {selectedRun.id}</h2>
            
            <div className="mb-4">
              <div className="flex items-center gap-2 mb-2">
                {getStatusIcon(selectedRun.status)}
                <span className="font-medium">Status: {selectedRun.status}</span>
              </div>
            </div>

            <div className="space-y-4">
              {selectedRun.job_order && selectedRun.job_order.length > 0 ? (
                // Display jobs in the order specified by job_order
                selectedRun.job_order.map((jobName) => {
                  const job = selectedRun.jobs[jobName];
                  if (!job) return null;
                  
                  const duration = job.started_at && job.ended_at 
                    ? Math.round((new Date(job.ended_at) - new Date(job.started_at)) / 1000)
                    : null;

                  return (
                    <div key={jobName} className="border rounded-lg p-4">
                      <div className="flex items-center gap-2 mb-3">
                        {getStatusIcon(job.status)}
                        <h3 className="font-semibold">{jobName}</h3>
                        <span className={`ml-auto px-2 py-1 rounded text-xs ${getStatusColor(job.status)}`}>
                          {job.status}
                        </span>
                      </div>
                      
                      {(job.started_at || duration) && (
                        <div className="text-xs text-gray-500 mb-3">
                          {job.started_at && (
                            <div>Started: {new Date(job.started_at).toLocaleString()}</div>
                          )}
                          {job.ended_at && (
                            <div>Ended: {new Date(job.ended_at).toLocaleString()}</div>
                          )}
                          {duration !== null && (
                            <div className="font-medium">Duration: {duration}s</div>
                          )}
                        </div>
                      )}
                      
                      {job.steps && job.steps.length > 0 && (
                        <div className="mb-3">
                          <p className="text-sm font-medium text-gray-700 mb-2">Steps:</p>
                          <ul className="space-y-1">
                            {job.steps.map((step, idx) => (
                              <li key={idx} className="text-sm text-gray-600 pl-4">
                                • {step.name}
                              </li>
                            ))}
                          </ul>
                        </div>
                      )}
                      
                      {job.output && (
                        <div>
                          <p className="text-sm font-medium text-gray-700 mb-2">Output:</p>
                          <pre className="bg-gray-900 text-green-400 p-3 rounded text-xs overflow-x-auto">
                            {job.output}
                          </pre>
                        </div>
                      )}
                    </div>
                  );
                })
              ) : (
                // Fallback to unordered display if job_order is not available
                Object.entries(selectedRun.jobs || {}).map(([jobName, job]) => {
                  const duration = job.started_at && job.ended_at 
                    ? Math.round((new Date(job.ended_at) - new Date(job.started_at)) / 1000)
                    : null;

                  return (
                    <div key={jobName} className="border rounded-lg p-4">
                      <div className="flex items-center gap-2 mb-3">
                        {getStatusIcon(job.status)}
                        <h3 className="font-semibold">{jobName}</h3>
                        <span className={`ml-auto px-2 py-1 rounded text-xs ${getStatusColor(job.status)}`}>
                          {job.status}
                        </span>
                      </div>
                      
                      {(job.started_at || duration) && (
                        <div className="text-xs text-gray-500 mb-3">
                          {job.started_at && (
                            <div>Started: {new Date(job.started_at).toLocaleString()}</div>
                          )}
                          {job.ended_at && (
                            <div>Ended: {new Date(job.ended_at).toLocaleString()}</div>
                          )}
                          {duration !== null && (
                            <div className="font-medium">Duration: {duration}s</div>
                          )}
                        </div>
                      )}
                      
                      {job.steps && job.steps.length > 0 && (
                        <div className="mb-3">
                          <p className="text-sm font-medium text-gray-700 mb-2">Steps:</p>
                          <ul className="space-y-1">
                            {job.steps.map((step, idx) => (
                              <li key={idx} className="text-sm text-gray-600 pl-4">
                                • {step.name}
                              </li>
                            ))}
                          </ul>
                        </div>
                      )}
                      
                      {job.output && (
                        <div>
                          <p className="text-sm font-medium text-gray-700 mb-2">Output:</p>
                          <pre className="bg-gray-900 text-green-400 p-3 rounded text-xs overflow-x-auto">
                            {job.output}
                          </pre>
                        </div>
                      )}
                    </div>
                  );
                })
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;