import { ApolloClient, InMemoryCache, HttpLink, ServerError, from } from '@apollo/client';
import { onError } from '@apollo/client/link/error';

const httpLink = new HttpLink({
  uri: '/query',
  credentials: 'include',
});

const errorLink = onError(({ graphQLErrors, networkError }) => {
  if ((networkError && (networkError as ServerError).statusCode === 401) || graphQLErrors?.some(error => error?.extensions?.code === 'UNAUTHENTICATED')) {
    window.location.href = '/';
  }
});

const client = new ApolloClient({
  link: from([errorLink, httpLink]),
  cache: new InMemoryCache(),
});

export default client;