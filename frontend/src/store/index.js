import { configureStore } from '@reduxjs/toolkit'
import authReducer from './authSlice'
import shareReducer from './shareSlice'

export const store = configureStore({
  reducer: {
    auth: authReducer,
    share: shareReducer
  }
})
