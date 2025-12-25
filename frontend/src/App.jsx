import React, { useState, useEffect, useCallback } from "react";
import { ArrowLeft } from "lucide-react";
import Header from "./components/Header";
import WorkflowList from "./components/WorkflowList";
import RunDetails from "./components/RunDetails";
import WorkflowHistory from "./components/WorkflowHistory";
import apiService from "./services/apiService";

export default function App() {
  const [workflows, setWorkflows] = useState([]);
  const [runs, setRuns] = useState([]);
  const [workflowStats, setWorkflowStats] = useState({});
  const [selectedRun, setSelectedRun] = useState(null);
  const [selectedWorkflow, setSelectedWorkflow] = useState(null);
  const [uploadStatus, setUploadStatus] = useState("");
  const [loading, setLoading] = useState(false);
  const [viewMode, setViewMode] = useState("workflows"); // workflows, history, runDetails

  // Fetch workflows
  const fetchWorkflows = useCallback(async () => {
    try {
      const data = await apiService.getWorkflows();
      setWorkflows(data || []);
    } catch (err) {
      console.error("Failed to fetch workflows:", err);
    }
  }, []);

  // Fetch workflow stats
  const fetchStats = useCallback(async (workflowNames) => {
    try {
      const stats = {};
      for (const name of workflowNames) {
        stats[name] = await apiService.getWorkflowStats(name);
      }
      setWorkflowStats(stats);
    } catch (err) {
      console.error("Failed to fetch stats:", err);
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

  // Handle workflow deletion
  const handleDeleteWorkflow = async (workflowName) => {
    if (!window.confirm(`Are you sure you want to delete "${workflowName}"? This cannot be undone.`)) {
      return;
    }

    try {
      await apiService.deleteWorkflow(workflowName);
      setUploadStatus(`✓ Workflow "${workflowName}" deleted successfully!`);
      await fetchWorkflows();
      setTimeout(() => setUploadStatus(""), 3000);
    } catch (err) {
      setUploadStatus(`✗ Failed to delete: ${err.message}`);
    }
  };

  // Handle run selection
  const handleSelectRun = (runId) => {
    fetchRunDetails(runId);
    setViewMode("runDetails");
  };

  // Handle view history
  const handleViewHistory = (workflowName) => {
    setSelectedWorkflow(workflowName);
    setViewMode("history");
  };

  // Handle back from history
  const handleBackFromHistory = () => {
    setViewMode("workflows");
    setSelectedWorkflow(null);
  };

  // Auto-refresh runs and stats
  useEffect(() => {
    fetchWorkflows();
    fetchRuns();

    const interval = setInterval(fetchRuns, 5000);
    return () => clearInterval(interval);
  }, [fetchWorkflows, fetchRuns]);

  // Fetch stats when workflows change
  useEffect(() => {
    if (workflows.length > 0) {
      fetchStats(workflows.map((w) => w.name));
    }
  }, [workflows, fetchStats]);

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

      <main className="max-w-7xl mx-auto px-4 py-8">
        {viewMode === "runDetails" && selectedRun ? (
          // Run Details View
          <div>
            <button
              onClick={() => {
                setViewMode("workflows");
                setSelectedRun(null);
              }}
              className="inline-flex items-center text-sm text-blue-600 hover:text-blue-800 mb-6"
            >
              <ArrowLeft className="w-4 h-4 mr-1" />
              Back to workflows
            </button>
            <RunDetails run={selectedRun} />
          </div>
        ) : viewMode === "history" && selectedWorkflow ? (
          // Workflow History View
          <WorkflowHistory
            workflowName={selectedWorkflow}
            onBack={handleBackFromHistory}
            onSelectRun={handleSelectRun}
          />
        ) : (
          // Main Dashboard View
          <div className="space-y-8">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">
                Workflows
              </h1>
              <p className="text-gray-600">
                Manage and monitor your CI/CD workflows
              </p>
            </div>

            <WorkflowList
              workflows={workflows}
              runs={runs}
              workflowStats={workflowStats}
              onUpload={handleUpload}
              onTrigger={handleTrigger}
              onSelectRun={handleSelectRun}
              onViewHistory={handleViewHistory}
              onDeleteWorkflow={handleDeleteWorkflow}
              uploadStatus={uploadStatus}
              loading={loading}
            />
          </div>
        )}
      </main>
    </div>
  );
}