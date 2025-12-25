import React, { useState, useEffect, useCallback } from "react";
import { ArrowLeft } from "lucide-react";
import Header from "./components/Header";
import WorkflowList from "./components/WorkflowList";
import RunList from "./components/RunList";
import RunDetails from "./components/RunDetails";
import apiService from "./services/apiService";

export default function App() {
  const [workflows, setWorkflows] = useState([]);
  const [runs, setRuns] = useState([]);
  const [selectedRun, setSelectedRun] = useState(null);
  const [uploadStatus, setUploadStatus] = useState("");
  const [loading, setLoading] = useState(false);

  // Fetch workflows
  const fetchWorkflows = useCallback(async () => {
    try {
      const data = await apiService.getWorkflows();
      setWorkflows(data || []);
    } catch (err) {
      console.error("Failed to fetch workflows:", err);
    }
  }, []);

  // Fetch runs
  const fetchRuns = useCallback(async () => {
    try {
      const data = await apiService.getRuns();
      setRuns(data || []);
    } catch (err) {
      console.error("Failed to fetch runs:", err);
    }
  }, []);

  // Fetch specific run details
  const fetchRunDetails = useCallback(async (runId) => {
    try {
      const data = await apiService.getRun(runId);
      setSelectedRun(data);
    } catch (err) {
      console.error("Failed to fetch run details:", err);
    }
  }, []);

  // Handle workflow upload
  const handleUpload = async (e) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      setUploadStatus("Uploading...");
      await apiService.uploadWorkflow(file);
      setUploadStatus("✓ Workflow uploaded successfully!");
      await fetchWorkflows();
      setTimeout(() => setUploadStatus(""), 3000);
    } catch (err) {
      setUploadStatus(`✗ Failed to upload: ${err.message}`);
    }
    e.target.value = "";
  };

  // Handle workflow trigger
  const handleTrigger = async (workflowName) => {
    try {
      setLoading(true);
      await apiService.triggerWorkflow(workflowName);
      await fetchRuns();
    } catch (err) {
      console.error("Failed to trigger workflow:", err);
      alert(`Failed to trigger workflow: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Handle run selection
  const handleSelectRun = (runId) => {
    fetchRunDetails(runId);
  };

  // Auto-refresh runs
  useEffect(() => {
    fetchWorkflows();
    fetchRuns();

    const interval = setInterval(fetchRuns, 5000);
    return () => clearInterval(interval);
  }, [fetchWorkflows, fetchRuns]);

  // Auto-refresh selected run details
  useEffect(() => {
    if (!selectedRun) return;

    const interval = setInterval(() => {
      fetchRunDetails(selectedRun.id);
    }, 3000);

    return () => clearInterval(interval);
  }, [selectedRun, fetchRunDetails]);

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <main className="max-w-7xl mx-auto px-4 py-6">
        {selectedRun ? (
          // Run Details View
          <div>
            <button
              onClick={() => setSelectedRun(null)}
              className="inline-flex items-center text-sm text-blue-600 hover:text-blue-800 mb-4"
            >
              <ArrowLeft className="w-4 h-4 mr-1" />
              Back to workflow runs
            </button>
            <RunDetails run={selectedRun} />
          </div>
        ) : (
          // Main Dashboard View
          <div className="space-y-6">
            <div>
              <h1 className="text-2xl font-semibold text-gray-900 mb-1">
                
              </h1>
              <p className="text-gray-600">
                Automate your workflow with CI/CD pipelines
              </p>
            </div>

            <WorkflowList
              workflows={workflows}
              onUpload={handleUpload}
              onTrigger={handleTrigger}
              uploadStatus={uploadStatus}
              loading={loading}
            />

            <RunList
              runs={runs}
              onSelectRun={handleSelectRun}
              onRefresh={fetchRuns}
            />
          </div>
        )}
      </main>
    </div>
  );
}