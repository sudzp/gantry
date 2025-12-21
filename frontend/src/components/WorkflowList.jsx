import React from 'react';
import { Play, Upload } from 'lucide-react';

export default function WorkflowList({ 
  workflows, 
  onUpload, 
  onTrigger, 
  uploadStatus, 
  loading 
}) {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-xl font-semibold mb-4">Workflows</h2>
      
      <div className="mb-4">
        <label className="flex items-center justify-center px-4 py-3 bg-blue-500 text-white rounded-lg cursor-pointer hover:bg-blue-600 transition">
          <Upload className="w-5 h-5 mr-2" />
          Upload Workflow YAML
          <input
            type="file"
            accept=".yml,.yaml"
            onChange={onUpload}
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
          <p className="text-gray-500 text-sm">
            No workflows yet. Upload one to get started!
          </p>
        ) : (
          workflows.map((wf) => (
            <div 
              key={wf.name} 
              className="flex items-center justify-between p-3 bg-gray-50 rounded border hover:bg-gray-100 transition"
            >
              <div className="flex-1">
                <span className="font-medium text-sm">{wf.name}</span>
                <div className="text-xs text-gray-500 mt-1">
                  {Object.keys(wf.jobs || {}).length} jobs
                </div>
              </div>
              <button
                onClick={() => onTrigger(wf.name)}
                disabled={loading}
                className="p-2 text-blue-600 hover:bg-blue-50 rounded disabled:opacity-50 transition"
                title="Trigger workflow"
              >
                <Play className="w-4 h-4" />
              </button>
            </div>
          ))
        )}
      </div>
    </div>
  );
}