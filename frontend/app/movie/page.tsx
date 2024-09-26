"use client"

import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import { useQuery, useMutation } from '@apollo/client';
import { DELETE_RATING, GET_MOVIE_BY_TITLE, RATE_MOVIE } from '@/graphql/queries';
import Link from 'next/link';
import getWikipediaImage from '@/lib/getImage';
import { Card, CardContent, Typography, Box, Rating, Button } from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';

export default function MovieDetail() {
  const [imageUrl, setImageUrl] = useState<string | null>(null);
  const [rating, setRating] = useState<number | null>(0);
  const searchParams = useSearchParams();
  const movieTitle = searchParams.get('title')?.split('/').pop()?.replace(/-/g, ' ');

  const { loading, error, data } = useQuery(GET_MOVIE_BY_TITLE, {
    variables: { title: movieTitle },
  });

  const [rateMovie] = useMutation(RATE_MOVIE, {
    onError: (error) => {
      console.error('Error rating movie:', error);
    },
  });

  const [deleteRating] = useMutation(DELETE_RATING, {
    onCompleted: () => {
      setRating(0);
    },
    onError: (error) => {
      console.error('Error deleting rating:', error);
    },
  });

  // Load the movie's image
  useEffect(() => {
    if (data?.movieByTitle?.wiki) {
      getWikipediaImage(data.movieByTitle.wiki).then(url => setImageUrl(url));
    }
  }, [data]);

  if (loading) return <Typography>Loading...</Typography>;
  if (error) return <Typography>Error: {error.message}</Typography>;

  const movie = data.movieByTitle;

  const handleRatingChange = (newValue: number | null) => {
    setRating(newValue);

    if (newValue) {
      rateMovie({ variables: { movieId: movie.id, score: newValue } });
    }
  };

  const handleDeleteRating = () => {
    deleteRating({ variables: { id: movie.ratingId } });
  };

  return (
    <Box sx={{ 
      display: 'flex', 
      flexDirection: 'column', 
      minHeight: '100vh', 
      position: 'relative', 
      padding: 2 
    }}>
      <Button
        component={Link}
        href="/"
        startIcon={<ArrowBackIcon />}
        sx={{ position: 'absolute', top: 16, left: 16 }}
        className='bg-sky-500 text-white'
      >
        Back to Movies
      </Button>
      <Box sx={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        flexGrow: 1 
      }}>
        <Card sx={{ 
          display: 'flex', 
          backgroundColor: 'white', 
          color: 'skyblue',
          maxWidth: '70%',
          width: '100%',
          borderRadius: 2
        }}>
          <Box sx={{ width: '40%', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <img src={imageUrl || ''} alt={movie.title} style={{ maxWidth: '100%', maxHeight: '300px', objectFit: 'cover' }} />
          </Box>
          <CardContent sx={{ width: '60%' }}>
            <Typography variant="h4" component="div" gutterBottom>
              {movie.title}
            </Typography>
            <Typography variant="body1" gutterBottom>
              <strong>Genre:</strong> {movie.genre}
            </Typography>
            <Typography variant="body1" gutterBottom>
              <strong>Year:</strong> {movie.year}
            </Typography>
            <Box sx={{ mt: 2 }}>
              <Typography component="legend"><strong>Rating</strong></Typography>
              <Button 
                onClick={() => handleRatingChange(0)} 
                sx={{ mt: 1, ml: 2 }} 
                color="secondary"
              >
                Dislike
              </Button>
              <Rating
                name="half-rating"
                value={rating}
                precision={0.5}
                onChange={(event, newValue) => handleRatingChange(newValue)}
                sx={{
                  '& .MuiRating-iconFilled': {
                    color: 'skyblue',
                  },
                  '& .MuiRating-iconHover': {
                    color: 'deepskyblue',
                  },
                }}
              />
              <Button 
                onClick={handleDeleteRating} 
                sx={{ mt: 1 }} 
                color="error"
              >
                Clear Rating
              </Button>
            </Box>
            <Typography variant="body1" gutterBottom className='line-clamp-4'>
              <strong>Plot:</strong> {movie.plot}
            </Typography>
          </CardContent>
        </Card>
      </Box>
    </Box>
  );
}
