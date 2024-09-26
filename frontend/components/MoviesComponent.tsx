import { useQuery } from '@apollo/client';
import { GET_MOVIES, SEARCH_MOVIES } from '../graphql/queries';
import { useEffect, useState } from 'react';
import MovieCard from './MovieCard';
import getWikipediaImage from '@/lib/getImage';
import { TextField, InputAdornment } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';

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
}

interface MovieImages {
    [key: string]: string;
}

export default function MoviesComponent() {
    const [page, setPage] = useState(1);
    const pageSize = 20;
    const [movieImages, setMovieImages] = useState<MovieImages>({});
    const [searchTerm, setSearchTerm] = useState('');

    const { loading, error, data, refetch } = useQuery<MoviesData>(
        searchTerm ? SEARCH_MOVIES : GET_MOVIES,
        {
            variables: searchTerm
                ? { query: searchTerm, page, pageSize }
                : { page, pageSize },
            notifyOnNetworkStatusChange: true,
        }
    );

    useEffect(() => {
        const movies = data?.movies?.edges || data?.searchMovies?.edges || [];
        movies.forEach(async ({ node }) => {
            if (node.wiki) {
                const imageUrl = await getWikipediaImage(node.wiki);
                if (imageUrl) {
                    setMovieImages(prev => ({ ...prev, [node.id]: imageUrl }));
                }
            }
        });
    }, [data]);

    const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
        const value = event.target.value;
        setSearchTerm(value);
        setPage(1);
        refetch(value ? { query: value, page: 1, pageSize } : { page: 1, pageSize });
    };

    const handlePageChange = (newPage: number) => {
        setPage(newPage);
        refetch(searchTerm
            ? { query: searchTerm, page: newPage, pageSize }
            : { page: newPage, pageSize }
        );
    };

    if (error) return <p>Error: {error.message}</p>;

    const movies = data?.movies?.edges || data?.searchMovies?.edges || [];
    const pageInfo = data?.movies?.pageInfo || data?.searchMovies?.pageInfo;

    return (
        <div className="container mx-auto px-4">
            <h2 className="text-2xl text-sky-400 font-bold mb-4">Movies</h2>
            <TextField
                fullWidth
                variant="outlined"
                placeholder="Search movies..."
                value={searchTerm}
                onChange={handleSearch}
                className="mb-6"
                slotProps={{
                    input: {
                        startAdornment: (
                            <InputAdornment position="start">
                                <SearchIcon />
                            </InputAdornment>
                        ),
                        className: 'bg-white mb-4',
                    }
                }}
            />
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