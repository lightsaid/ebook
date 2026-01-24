

import { create } from "zustand";


interface AppState {
  isCollapsed: boolean;
  onMenuChange: () => void;
}

export const useAppStore = create<AppState>((set) => ({
  isCollapsed: false, // 默认不收起
  onMenuChange: () => set((state) => ({ isCollapsed: !state.isCollapsed })),
}));

