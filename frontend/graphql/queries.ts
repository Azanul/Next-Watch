import { gql } from '@apollo/client';

export const GET_MOVIES = gql`
  query GetMovies($page: Int!, $pageSize: Int!) {
    movies(page: $page, pageSize: $pageSize) {
      edges {
        node {
          id
          title
          genre
          year
          wiki
        }
      }
      pageInfo {
        hasNextPage
        hasPreviousPage
      }
      totalCount
    }
  }
`;

export const GET_RECOMMENDATIONS = gql`
  query GetRecommendations($page: Int!, $pageSize: Int!) {
    recommendations(page: $page, pageSize: $pageSize) {
      edges {
        node {
          id
          title
          genre
          year
        }
      }
      pageInfo {
        hasNextPage
        hasPreviousPage
      }
      totalCount
    }
  }
`;

export const GET_MOVIE_BY_TITLE = gql`
  query GetMovieByTitle($title: String!) {
    movieByTitle(title: $title) {
      id
      title
      genre
      year
      wiki
      plot
      cast
    }
  }
`;

export const SEARCH_MOVIES = gql`
  query SearchMovies($query: String!, $page: Int!, $pageSize: Int!) {
    searchMovies(query: $query, page: $page, pageSize: $pageSize) {
      edges {
        node {
          id
          title
          genre
          year
          wiki
        }
      }
      pageInfo {
        hasNextPage
        hasPreviousPage
      }
      totalCount
    }
  }
`;

export const RATE_MOVIE = gql`
  mutation RateMovie($movieId: ID!, $score: Float!) {
    rateMovie(movieId: $movieId, score: $score) {
      id
      score
    }
  }
`;

export const DELETE_RATING = gql`
  mutation DeleteRating($id: ID!) {
    deleteRating(id: $id)
  }
`;