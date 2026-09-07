import { ElLoading } from "element-plus";

let loading = null;

export function useLoading() {
  const show = (options = {}) => {
    if (loading) return;

    loading = ElLoading.service({
      lock: true,
      text: "Loading...",
      background: "rgba(255,255,255,0.2)",
      ...options,
    });
  };

  const hide = () => {
    loading?.close();
    loading = null;
  };

  return {
    show,
    hide,
  };
}