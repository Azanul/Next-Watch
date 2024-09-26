import { useQuery } from '@apollo/client';
import { GET_MOVIES, GET_RECOMMENDATIONS, SEARCH_MOVIES } from '../graphql/queries';
import { useEffect, useState } from 'react';
import MovieCard from './MovieCard';
import getWikipediaImage from '@/lib/getImage';

interface MovieNode {
    wiki: string;
    id: string;
    title: string;
    genre: string;
    year: number;
}

interface MoviesData {
    movies?: {
        edges: {
            node: MovieNode;
        }[];
        pageInfo: {
            hasNextPage: boolean;
            hasPreviousPage: boolean;
        };
        totalCount: number;
    };
    searchMovies?: {
        edges: {
            node: MovieNode;
        }[];
        pageInfo: {
            hasNextPage: boolean;
            hasPreviousPage: boolean;
        };
        totalCount: number;
    };
    recommendations?: {
        edges: {
            node: MovieNode;
        }[];
        pageInfo: {
            hasNextPage: boolean;
            hasPreviousPage: boolean;
        };
        totalCount: number;
    };
}

interface MovieImages {
    [key: string]: string;
}

interface MoviesComponentProps {
    queryType: 'GET_MOVIES' | 'SEARCH_MOVIES' | 'GET_RECOMMENDATIONS';
    searchTerm?: string;
}

export default function MoviesComponent({ queryType, searchTerm }: MoviesComponentProps) {
    const [page, setPage] = useState(1);
    const pageSize = 20;
    const [movieImages, setMovieImages] = useState<MovieImages>({});

    const { loading, error, data, refetch } = useQuery<MoviesData>(
        queryType === 'GET_RECOMMENDATIONS'
            ? GET_RECOMMENDATIONS
            : queryType === 'SEARCH_MOVIES' && searchTerm ? SEARCH_MOVIES : GET_MOVIES,
        {
            variables: searchTerm
                ? { query: searchTerm, page, pageSize }
                : { page, pageSize },
            notifyOnNetworkStatusChange: true,
            fetchPolicy: queryType === 'GET_RECOMMENDATIONS' ? 'no-cache' : 'cache-first',
        }
    );

    useEffect(() => {
        const movies = data?.movies?.edges || data?.searchMovies?.edges || data?.recommendations?.edges || [];
        movies.forEach(async ({ node }) => {
            if (node.wiki) {
                const imageUrl = await getWikipediaImage(node.wiki);
                if (imageUrl) {
                    setMovieImages(prev => ({ ...prev, [node.id]: imageUrl }));
                }
            }
        });
    }, [data]);

    const handlePageChange = (newPage: number) => {
        setPage(newPage);
        refetch(searchTerm
            ? { query: searchTerm, page: newPage, pageSize }
            : { page: newPage, pageSize }
        );
    };

    if (error) return <p>Error: {error.message}</p>;

    const movies = data?.movies?.edges || data?.searchMovies?.edges || data?.recommendations?.edges || [];
    const pageInfo = data?.movies?.pageInfo || data?.searchMovies?.pageInfo || data?.recommendations?.pageInfo;

    return (
        <div className="container mx-auto px-4">
            {loading && <p>Loading...</p>}
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                {movies.map(({ node }) => (
                    <MovieCard
                        key={node.id}
                        id={node.id}
                        title={node.title}
                        genre={node.genre}
                        year={node.year}
                        imageUrl={movieImages[node.id] || ''}
                    />
                ))}
            </div>
            <div className="mt-8 flex justify-between">
                <button
                    onClick={() => handlePageChange(Math.max(1, page - 1))}
                    disabled={!pageInfo?.hasPreviousPage || loading}
                    className="bg-sky-500 text-white px-4 py-2 rounded disabled:bg-gray-300"
                >
                    Previous
                </button>
                <button
                    onClick={() => handlePageChange(page + 1)}
                    disabled={!pageInfo?.hasNextPage || loading}
                    className="bg-sky-500 text-white px-4 py-2 rounded disabled:bg-gray-300"
                >
                    Next
                </button>
            </div>
        </div>
    );
}
