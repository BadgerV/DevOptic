import { configureStore, combineReducers } from "@reduxjs/toolkit";
import { setupListeners } from "@reduxjs/toolkit/query";
import storage from "redux-persist/lib/storage";
import { persistStore, persistReducer } from "redux-persist";

// Import domain APIs
import { endpointMonitoringApi } from "@domains/endpoint-monitoring/store/endpointMonitoringApi";
import endpointMonitoringreducer from "@domains/endpoint-monitoring/store/endpointMonitoringSlice";
import { gitlabApi, gitlabReducer } from "@/domains/gitlab/store";
import { authApi, authReducer } from "../../shared/services/api/authApi";
import { rbacApi } from "@/domains/rbac/store/rbacApi";

/* -------------------- Combine reducers -------------------- */
const rootReducer = combineReducers({
  // API reducers
  [endpointMonitoringApi.reducerPath]: endpointMonitoringApi.reducer,
  [authApi.reducerPath]: authApi.reducer,
  [rbacApi.reducerPath]: rbacApi.reducer,
  [gitlabApi.reducerPath]: gitlabApi.reducer,

  // Regular slices
  endpointMonitoring: endpointMonitoringreducer,
  auth: authReducer,
  gitlab: gitlabReducer,
});

/* -------------------- Persist config -------------------- */
const persistConfig = {
  key: "root",
  storage,
  whitelist: ["auth"], // only persist auth slice
};

const persistedReducer = persistReducer(persistConfig, rootReducer);

/* -------------------- Store -------------------- */
export const store = configureStore({
  reducer: persistedReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: ["persist/PERSIST", "persist/REHYDRATE"],
      },
    }).concat(
      endpointMonitoringApi.middleware,
      authApi.middleware,
      rbacApi.middleware,
      gitlabApi.middleware
    ),
});

export const persistor = persistStore(store);

// Enable RTK Query refetchOnFocus/refetchOnReconnect
setupListeners(store.dispatch);

/* -------------------- Types -------------------- */
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
