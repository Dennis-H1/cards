import { useCallback, useEffect, useState } from "react";
import { Route, Routes } from "react-router-dom";
import { listActivity } from "./api/client";
import { BottomNav } from "./components/BottomNav";
import { ActivityPage } from "./pages/ActivityPage";
import { BrowsePage } from "./pages/BrowsePage";
import { CardDetailPage } from "./pages/CardDetailPage";
import { ReviewPage } from "./pages/ReviewPage";
import { TagOverviewPage } from "./pages/TagOverviewPage";

function App() {
  const [unseenCount, setUnseenCount] = useState(0);

  const refreshUnseenCount = useCallback(() => {
    listActivity(true)
      .then((events) => setUnseenCount(events.length))
      .catch(() => {
        // Activity badge is a nice-to-have; ignore transient fetch errors.
      });
  }, []);

  useEffect(() => {
    refreshUnseenCount();
  }, [refreshUnseenCount]);

  return (
    <div className="app-layout">
      <div className="app-content">
        <Routes>
          <Route path="/" element={<ReviewPage />} />
          <Route path="/browse" element={<BrowsePage />} />
          <Route path="/cards/:id" element={<CardDetailPage />} />
          <Route path="/tags/:name" element={<TagOverviewPage />} />
          <Route path="/activity" element={<ActivityPage onSeen={refreshUnseenCount} />} />
        </Routes>
      </div>
      <BottomNav unseenCount={unseenCount} />
    </div>
  );
}

export default App;
