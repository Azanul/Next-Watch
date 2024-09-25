"use client";

import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function MovieDetail({ movie }: { movie: any }) {
  const router = useRouter();
  const [rating, setRating] = useState(0);

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-4">{movie.title}</h1>
      <div className="mb-4">
        <img src={movie.imageUrl} alt={movie.title} className="w-full max-w-md mx-auto" />
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
      <button
        onClick={() => router.back()}
        className="bg-sky-500 text-white px-4 py-2 rounded"
      >
        Back to Movies
      </button>
    </div>
  );
}
