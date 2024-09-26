import { useState } from "react";
import MoviesComponent from "./MoviesComponent";
import { InputAdornment, TextField } from "@mui/material";
import SearchIcon from '@mui/icons-material/Search';

export default function AllMoviesScreen() {
    const [searchTerm, setSearchTerm] = useState('');

    const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
        const value = event.target.value;
        setSearchTerm(value);
    };

    return (
        <div className="container mx-auto px-4">
            <TextField
                fullWidth
                variant="outlined"
                placeholder="Search movies..."
                value={searchTerm}
                onChange={handleSearch}
                className="mb-6"
                InputProps={{
                    startAdornment: (
                        <InputAdornment position="start">
                            <SearchIcon />
                        </InputAdornment>
                    ),
                }}
            />
            <MoviesComponent queryType={searchTerm ? 'SEARCH_MOVIES' : 'GET_MOVIES'} searchTerm={searchTerm}></MoviesComponent>
        </div>
    )
}