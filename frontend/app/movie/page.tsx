"use client"

import { useSearchParams } from 'next/navigation'
import { useQuery } from '@apollo/client';
import { GET_MOVIE_BY_TITLE } from '@/graphql/queries';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import getWikipediaImage from '@/lib/getImage';


export default function MovieDetail() {
  const [imageUrl, setImageUrl] = useState<string | null>(null);
  const [rating, setRating] = useState(0);
  const searchParams = useSearchParams();
  const movieTitle = searchParams.get('title')?.split('/').pop()?.replace(/-/g, ' ');

  const { loading, error, data } = useQuery(GET_MOVIE_BY_TITLE, {
    variables: { title: movieTitle },
  });

  useEffect(() => {
    if (data?.movieByTitle?.wiki) {
      getWikipediaImage(data.movieByTitle.wiki).then(url => setImageUrl(url));
    }
  }, [data]);


  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error.message}</p>;

  const movie = data.movieByTitle;

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-4">{movie.title}</h1>
      <div className="mb-4">
        <img src={imageUrl || ''} alt={movie.title} className="w-full max-w-md mx-auto" />
      </div>
      <p className="mb-2"><strong>Genre:</strong> {movie.genre}</p>
      <p className="mb-2"><strong>Year:</strong> {movie.year}</p>
      <div className="mb-4">
        <h2 className="text-xl font-semibold mb-2">Rating</h2>
        <div className="flex items-center">
          {[1, 2, 3, 4, 5].map((star) => (
            <button
              key={star}
              onClick={() => setRating(star)}
              className={`text-3xl ${star <= rating ? 'text-yellow-400' : 'text-gray-300'}`}
            >
              ★
            </button>
          ))}
        </div>
      </div>
      <div
        className="bg-sky-500 text-white px-4 py-2 rounded"
      >
        <Link href={"/"}>Back to Movies</Link>
      </div>
    </div>
  );
}
