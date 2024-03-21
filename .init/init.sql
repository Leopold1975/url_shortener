SELECT 'CREATE DATABASE urls' 
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'urls')\gexec