"use client"

import { useQuery } from '@apollo/client';
import { GET_MOVIES } from '../graphql/queries';
import { useEffect, useState } from 'react';

interface MovieNode {
    wiki: string;
    id: string;
    title: string;
    genre: string;
    year: number;
}

interface MoviesData {
    movies: {
        edges: {
            node: MovieNode;
        }[];
        pageInfo: {
            hasNextPage: boolean;
            hasPreviousPage: boolean;
        };
    };
}

interface MovieImages {
    [key: string]: string;
}

async function getWikipediaImage(wikiUrl: string): Promise<string | null> {
    const pageTitle = wikiUrl.split('/wiki/')[1];
    const apiUrl = `https://en.wikipedia.org/w/api.php?action=parse&page=${pageTitle}&format=json&prop=text&origin=*`;

    try {
        const response = await fetch(apiUrl);
        const data = await response.json();
        const pageContent = data.parse.text['*'];
        const parser = new DOMParser();
        const doc = parser.parseFromString(pageContent, 'text/html');
        const infobox = doc.querySelector('.infobox');
        const img = infobox?.querySelector('img');
        return img ? `https:${img.getAttribute('src')}` : null;
    } catch (error) {
        console.error('Error fetching Wikipedia image:', error);
        return null;
    }
}

export default function MoviesComponent() {
    const [page, setPage] = useState(1);
    const pageSize = 20;
    const [movieImages, setMovieImages] = useState<MovieImages>({});

    const { loading, error, data } = useQuery<MoviesData>(GET_MOVIES, {
        variables: { page, pageSize },
    });

    useEffect(() => {
        if (data?.movies?.edges) {
            data.movies.edges.forEach(async ({ node }) => {
                if (node.wiki) {
                    const imageUrl = await getWikipediaImage(node.wiki);
                    if (imageUrl) {
                        setMovieImages(prev => ({ ...prev, [node.id]: imageUrl }));
                    }
                }
            });
        }
    }, [data]);

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error: {error.message}</p>;

    const movies = data?.movies?.edges || [];
    const pageInfo = data?.movies?.pageInfo;

    return (
        <div className="container mx-auto px-4">
            <h2 className="text-2xl text-sky-400 font-bold mb-4">Movies</h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                {movies.map(({ node }) => (
                    <div key={node.id} className="rounded-lg shadow-md overflow-hidden">
                        <div
                            className="h-48 bg-cover bg-center"
                            style={{ backgroundImage: `url(${movieImages[node.id] || ''})` }}

                        />
                        <div key={node.id} className="bg-sky-400 p-2 h-full rounded shadow-md">
                            <h3 className="font-semibold text-white">{node.title}</h3>
                            <p className="text-sky-100">{node.genre} - {node.year}</p>
                        </div>
                    </div>
                ))}
            </div>
            <div className="mt-8 flex justify-between">
                <button
                    onClick={() => setPage(p => Math.max(1, p - 1))}
                    disabled={!pageInfo?.hasPreviousPage}
                    className="bg-sky-500 text-white px-4 py-2 rounded disabled:bg-gray-300"
                >
                    Previous
                </button>
                <button
                    onClick={() => setPage(p => p + 1)}
                    disabled={!pageInfo?.hasNextPage}
                    className="bg-sky-500 text-white px-4 py-2 rounded disabled:bg-gray-300"
                >
                    Next
                </button>
            </div>
        </div>
    );
}
