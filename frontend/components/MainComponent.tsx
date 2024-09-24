"use client"

import { useState } from 'react';
import MoviesComponent from './MoviesComponent';
import RecommendationsComponent from './RecommendationsComponent';

export default function MainComponent() {
    const [activeTab, setActiveTab] = useState('movies');

    return (
        <div className="w-full max-w-4xl mx-auto">
            <div className="mb-4 justify-center flex">
                <div className='inline-flex rounded-md shadow-sm' role="group">
                    <button
                        type="button"
                        className={`px-4 py-2 text-sm font-medium rounded-l-lg ${
                            activeTab === 'movies'
                                ? 'bg-sky-500 text-white'
                                : 'bg-white text-gray-700 hover:bg-gray-50'
                        } border border-gray-200`}
                        onClick={() => setActiveTab('movies')}
                    >
                        Movies
                    </button>
                    <button
                        type="button"
                        className={`px-4 py-2 text-sm font-medium rounded-r-lg ${
                            activeTab === 'recommendations'
                                ? 'bg-blue-500 text-white'
                                : 'bg-white text-gray-700 hover:bg-gray-50'
                        } border border-gray-200`}
                        onClick={() => setActiveTab('recommendations')}
                    >
                        Recommendations
                    </button>
                </div>
            </div>
            <div className="bg-white p-4 rounded-lg shadow">
                {activeTab === 'movies' ? <MoviesComponent /> : <RecommendationsComponent />}
            </div>
        </div>
    );
}