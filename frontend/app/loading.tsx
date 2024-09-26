import { CircularProgress, Typography } from '@mui/material';

export default function Loading() {
  return (
    <div className="flex flex-col items-center justify-center h-screen bg-gradient-to-b from-sky-400 to-sky-200">
      <CircularProgress
        size={60}
        thickness={4}
        sx={{
          color: 'white',
          marginBottom: '1rem',
        }}
      />
      <Typography
        variant="h5"
        component="h1"
        className="text-white font-semibold animate-pulse"
      >
        Loading...
      </Typography>
    </div>
  );
};