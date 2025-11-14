import { createSlice } from '@reduxjs/toolkit'

const initialState = {
  shares: [],
  currentShare: null,
  loading: false
}

const shareSlice = createSlice({
  name: 'share',
  initialState,
  reducers: {
    setShares: (state, action) => {
      state.shares = action.payload
    },
    setCurrentShare: (state, action) => {
      state.currentShare = action.payload
    },
    setLoading: (state, action) => {
      state.loading = action.payload
    }
  }
})

export const { setShares, setCurrentShare, setLoading } = shareSlice.actions
export default shareSlice.reducer
