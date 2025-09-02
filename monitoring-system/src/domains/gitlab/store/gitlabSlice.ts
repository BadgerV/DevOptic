import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface GitlabState {
  selectedPipelineUnitId: string | null;
  isStatusModalOpen: boolean;
  statusModalPipelineRunId: string | null;
  approvalFilter: 'all' | 'pending' | 'approved' | 'rejected';
  historyStatusFilter: 'all' | 'running' | 'success' | 'failed' | 'pending';
  searchTerm: string;
  ui: {
    isSidebarOpen: boolean;
  };
}

const initialState: GitlabState = {
  selectedPipelineUnitId: null,
  isStatusModalOpen: false,
  statusModalPipelineRunId: null,
  approvalFilter: 'pending',
  historyStatusFilter: 'all',
  searchTerm: '',
  ui: {
    isSidebarOpen: true,
  },
};

const gitlabSlice = createSlice({
  name: 'gitlab',
  initialState,
  reducers: {
    setSelectedPipelineUnitId(state, action: PayloadAction<string | null>) {
      state.selectedPipelineUnitId = action.payload;
    },
    openStatusModal(state, action: PayloadAction<string>) {
      state.isStatusModalOpen = true;
      state.statusModalPipelineRunId = action.payload;
    },
    closeStatusModal(state) {
      state.isStatusModalOpen = false;
      state.statusModalPipelineRunId = null;
    },
    setApprovalFilter(state, action: PayloadAction<'all' | 'pending' | 'approved' | 'rejected'>) {
      state.approvalFilter = action.payload;
    },
    setHistoryStatusFilter(state, action: PayloadAction<'all' | 'running' | 'success' | 'failed' | 'pending'>) {
      state.historyStatusFilter = action.payload;
    },
    setSearchTerm(state, action: PayloadAction<string>) {
      state.searchTerm = action.payload;
    },
    toggleSidebar(state) {
      state.ui.isSidebarOpen = !state.ui.isSidebarOpen;
    },
  },
});

export const {
  setSelectedPipelineUnitId,
  openStatusModal,
  closeStatusModal,
  setApprovalFilter,
  setHistoryStatusFilter,
  setSearchTerm,
  toggleSidebar,
} = gitlabSlice.actions;

export default gitlabSlice.reducer;