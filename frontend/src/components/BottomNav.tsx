import { NavLink } from "react-router-dom";

export function BottomNav({ unseenCount }: { unseenCount: number }) {
  return (
    <nav className="bottom-nav">
      <NavLink to="/" end className={({ isActive }) => (isActive ? "active" : "")}>
        Review
      </NavLink>
      <NavLink to="/browse" className={({ isActive }) => (isActive ? "active" : "")}>
        Browse
      </NavLink>
      <NavLink to="/activity" className={({ isActive }) => (isActive ? "active" : "")}>
        Activity{unseenCount > 0 ? <span className="badge">{unseenCount}</span> : null}
      </NavLink>
    </nav>
  );
}
