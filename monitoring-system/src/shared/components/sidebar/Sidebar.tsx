import React, { useState, useCallback, useMemo, useRef, useEffect } from "react";
import "./sidebar.css";
import { useAppDispatch, useAppSelector } from "@/shared/hooks/redux";
import { setFilterTypeByRunning } from "@/domains/endpoint-monitoring/store/endpointMonitoringSlice";
import { useNavigate, useLocation } from "react-router-dom";

type Environment = "QA" | "Demo";
type ServiceFilter = "All" | "Running" | "Stopped";

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({ isOpen, onClose }) => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const location = useLocation();
  
  // Refs for drag functionality
  const sidebarRef = useRef<HTMLElement>(null);
  const isDraggingRef = useRef(false);
  const startXRef = useRef(0);
  const startWidthRef = useRef(0);

  // Redux state - optimized selector
  const endpointData = useAppSelector(useCallback((state) => ({
    totalNumberOfEndpoints: state.endpointMonitoring.totalNumberOfEndpoints,
    runningEndpoints: state.endpointMonitoring.runningEndpoints,
  }), []));

  // Memoize calculated values
  const notRespondingEndpoints = useMemo(
    () => Math.max(endpointData.totalNumberOfEndpoints - endpointData.runningEndpoints, 0),
    [endpointData.totalNumberOfEndpoints, endpointData.runningEndpoints]
  );

  // Local state
  const [isDarkMode, setIsDarkMode] = useState<boolean>(true);
  const [environment, setEnvironment] = useState<Environment>("QA");
  const [activeServiceFilter, setActiveServiceFilter] = useState<ServiceFilter>("Stopped");
  const [sidebarWidth, setSidebarWidth] = useState<number>(280);

  // Drag functionality for resizing sidebar
  const handleMouseDown = useCallback((e: React.MouseEvent) => {
    if (!sidebarRef.current) return;
    
    isDraggingRef.current = true;
    startXRef.current = e.clientX;
    startWidthRef.current = sidebarRef.current.offsetWidth;
    
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
    document.body.style.cursor = 'col-resize';
    document.body.style.userSelect = 'none';
  }, []);

  const handleMouseMove = useCallback((e: MouseEvent) => {
    if (!isDraggingRef.current) return;
    
    requestAnimationFrame(() => {
      const deltaX = e.clientX - startXRef.current;
      const newWidth = Math.max(200, Math.min(500, startWidthRef.current + deltaX));
      setSidebarWidth(newWidth);
    });
  }, []);

  const handleMouseUp = useCallback(() => {
    isDraggingRef.current = false;
    document.removeEventListener('mousemove', handleMouseMove);
    document.removeEventListener('mouseup', handleMouseUp);
    document.body.style.cursor = '';
    document.body.style.userSelect = '';
  }, [handleMouseMove]);

  // Cleanup event listeners
  useEffect(() => {
    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
      document.body.style.cursor = '';
      document.body.style.userSelect = '';
    };
  }, [handleMouseMove, handleMouseUp]);

  // Optimized callbacks
  const setFilterByRunning = useCallback((type: string | boolean) => {
    dispatch(setFilterTypeByRunning(type));
  }, [dispatch]);

  const toggleTheme = useCallback((): void => {
    setIsDarkMode(prev => !prev);
    document.body.classList.toggle("light-mode");
  }, []);

  const handleServiceFilterClick = useCallback((filter: ServiceFilter): void => {
    if (filter !== activeServiceFilter) {
      setActiveServiceFilter(filter);
    }
  }, [activeServiceFilter]);

  // Navigation handlers
  const handleNavigate = useCallback((path: string) => {
    if (location.pathname !== path) {
      navigate(path);
    }
  }, [navigate, location.pathname]);

  // Memoized class calculations
  const getServiceFilterClass = useCallback((filter: ServiceFilter): string => {
    const baseClass = `filter-item filter-${filter.toLowerCase()}`;
    return filter === activeServiceFilter ? `${baseClass} active` : baseClass;
  }, [activeServiceFilter]);

  const getNavItemClass = useCallback((path: string): string => {
    return location.pathname === path ? "nav-item active" : "nav-item";
  }, [location.pathname]);

  // Filter handlers
  const handleAllFilter = useCallback((e: React.MouseEvent) => {
    e.preventDefault();
    handleServiceFilterClick("All");
    setFilterByRunning("all");
  }, [handleServiceFilterClick, setFilterByRunning]);

  const handleRunningFilter = useCallback((e: React.MouseEvent) => {
    e.preventDefault();
    handleServiceFilterClick("Running");
    setFilterByRunning(true);
  }, [handleServiceFilterClick, setFilterByRunning]);

  const handleStoppedFilter = useCallback((e: React.MouseEvent) => {
    e.preventDefault();
    handleServiceFilterClick("Stopped");
    setFilterByRunning(false);
  }, [handleServiceFilterClick, setFilterByRunning]);

  // Navigation handlers
  const navigateToDashboard = useCallback((e: React.MouseEvent) => {
    e.preventDefault();
    handleNavigate("/");
  }, [handleNavigate]);

  const navigateToSettings = useCallback((e: React.MouseEvent) => {
    e.preventDefault();
    handleNavigate("/settings");
  }, [handleNavigate]);

  const navigateToAbout = useCallback((e: React.MouseEvent) => {
    e.preventDefault();
    handleNavigate("/about");
  }, [handleNavigate]);

  // Close handler
  const handleClose = useCallback((e: React.MouseEvent) => {
    e.preventDefault();
    onClose();
  }, [onClose]);

  // Memoized service filters
  const serviceFilters = useMemo(() => [
    {
      key: "All",
      filter: "All" as ServiceFilter,
      text: "All",
      count: endpointData.totalNumberOfEndpoints,
      onClick: handleAllFilter,
      className: getServiceFilterClass("All")
    },
    {
      key: "Running",
      filter: "Running" as ServiceFilter,
      text: "Running",
      count: endpointData.runningEndpoints,
      onClick: handleRunningFilter,
      className: getServiceFilterClass("Running")
    },
    {
      key: "Stopped",
      filter: "Stopped" as ServiceFilter,
      text: "Not responding",
      count: notRespondingEndpoints,
      onClick: handleStoppedFilter,
      className: getServiceFilterClass("Stopped")
    }
  ], [
    endpointData.totalNumberOfEndpoints,
    endpointData.runningEndpoints,
    notRespondingEndpoints,
    handleAllFilter,
    handleRunningFilter,
    handleStoppedFilter,
    getServiceFilterClass
  ]);

  // Memoized navigation items
  const navigationItems = useMemo(() => [
    { 
      path: "/", 
      text: "Dashboard", 
      onClick: navigateToDashboard,
      className: getNavItemClass("/")
    },
    { 
      path: "/settings", 
      text: "Settings", 
      onClick: navigateToSettings,
      className: getNavItemClass("/settings")
    },
    { 
      path: "/about", 
      text: "About", 
      onClick: navigateToAbout,
      className: getNavItemClass("/about")
    }
  ], [
    navigateToDashboard, 
    navigateToSettings, 
    navigateToAbout,
    getNavItemClass
  ]);

  // Memoized sidebar style
  const sidebarStyle = useMemo(() => ({
    width: isOpen ? `${sidebarWidth}px` : '0px',
    transition: isDraggingRef.current ? 'none' : 'width 0.3s ease'
  }), [isOpen, sidebarWidth]);

  return (
    <aside 
      ref={sidebarRef}
      className={`sidebar ${isOpen ? "sidebar-open" : "sidebar-closed"}`}
      style={sidebarStyle}
    >
      {/* Resize Handle */}
      {isOpen && (
        <div 
          className="sidebar-resize-handle"
          onMouseDown={handleMouseDown}
          style={{
            position: 'absolute',
            right: 0,
            top: 0,
            bottom: 0,
            width: '4px',
            backgroundColor: 'transparent',
            cursor: 'col-resize',
            zIndex: 1000
          }}
        />
      )}

      {/* Logo Section */}
      <div className="logo-section">
        <div className="logo">
          <div className="logo-text">
            <span className="logo-title">Service Monitor</span>
            <span className="logo-subtitle">Dashboard v2.0</span>
          </div>
        </div>
      </div>

      {/* Navigation Section */}
      <div className="nav-section">
        <div className="section-header">
          <span className="section-title">Navigation</span>
        </div>
        <nav>
          <ul>
            {navigationItems.map(({ path, text, onClick, className }) => (
              <li
                key={path}
                className={className}
                onClick={onClick}
                role="button"
                tabIndex={0}
              >
                <span className="nav-text">{text}</span>
              </li>
            ))}
          </ul>
        </nav>
      </div>

      {/* Services Filter Section */}
      <div className="services-section">
        <div className="section-header">
          <span className="section-title">Services</span>
        </div>
        <div className="services-filter">
          {serviceFilters.map(({ key, text, count, onClick, className }) => (
            <div
              key={key}
              className={className}
              onClick={onClick}
              role="button"
              tabIndex={0}
            >
              <span className="filter-text">{text}</span>
              <span className="filter-count">{count}</span>
            </div>
          ))}
        </div>
      </div>

      {/* Bottom Section */}
      <div className="bottom-section">
        <div className="close-button-container">
          <button onClick={handleClose} className="close-btn">
            <span>Close</span>
          </button>
        </div>
      </div>
    </aside>
  );
};

// Memoize the component with custom comparison
export default React.memo(Sidebar, (prevProps, nextProps) => {
  return prevProps.isOpen === nextProps.isOpen && prevProps.onClose === nextProps.onClose;
});