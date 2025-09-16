# GitHub Actions Cache Management

This document describes the cache management strategies and tools for the GCP Automation API project to help manage GitHub Actions cache storage limits.

## Current Situation

The repository is approaching the 10 GB GitHub Actions cache storage limit (currently at 11.71 GB). GitHub automatically evicts least recently used caches, but this can cause build performance issues.

## Cache Management Solutions

### 1. Automated Cache Cleanup Workflow

**File**: `.github/workflows/cache-cleanup.yml`

This workflow runs daily at 2:00 AM UTC to automatically clean up old caches.

#### Features:
- **Conservative cleanup**: Deletes caches older than 7 days
- **Aggressive cleanup**: Deletes caches older than 3 days, then oldest until under size limit
- **Emergency cleanup**: Deletes ALL caches (manual trigger only)
- **Smart cleanup**: Removes duplicate Go module caches and large Docker buildx caches

#### Manual Triggering:
```bash
# Conservative cleanup (default)
gh workflow run cache-cleanup.yml

# Aggressive cleanup to 7GB limit
gh workflow run cache-cleanup.yml -f cleanup_strategy=aggressive -f max_cache_size_gb=7

# Emergency cleanup (deletes everything)
gh workflow run cache-cleanup.yml -f cleanup_strategy=emergency
```

### 2. Manual Cache Cleanup Script

**File**: `scripts/cleanup-caches.sh`

A comprehensive script for manual cache management with interactive options.

#### Usage:
```bash
# Interactive mode
./scripts/cleanup-caches.sh

# Automated aggressive cleanup
./scripts/cleanup-caches.sh stuartshay/gcp-automation-api 7 aggressive

# List caches only
./scripts/cleanup-caches.sh stuartshay/gcp-automation-api 8 list
```

#### Prerequisites:
- GitHub CLI (`gh`) installed and authenticated
- `jq` for JSON processing
- `bc` for calculations

### 3. Optimized Cache Configurations

#### Go Module Caching Improvements:
- **Consistent cache keys**: Using `go.mod` and `go.sum` instead of wildcard patterns
- **Version-specific keys**: Including Go version in cache keys to prevent conflicts
- **Hierarchical restore keys**: Better cache hit rates with fallback options

**Before**:
```yaml
key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

**After**:
```yaml
key: ${{ runner.os }}-go-${{ env.GO_VERSION }}-${{ hashFiles('go.mod', 'go.sum') }}
restore-keys: |
  ${{ runner.os }}-go-${{ env.GO_VERSION }}-
  ${{ runner.os }}-go-
```

#### Docker Build Cache Optimization:
- **Changed from `mode=max` to `mode=min`**: Reduces cache size by only storing essential layers
- **Removed inline caching**: Prevents redundant cache storage

**Before**:
```yaml
cache-to: type=gha,mode=max
```

**After**:
```yaml
cache-to: type=gha,mode=min
cache-from-inline: false
```

### 4. Cache Size Monitoring

Added cache size monitoring to workflows to track usage:

```yaml
- name: Check cache sizes
  run: |
    echo "Go build cache size:"
    du -sh ~/.cache/go-build 2>/dev/null || echo "No go-build cache found"
    echo "Go module cache size:"
    du -sh ~/go/pkg/mod 2>/dev/null || echo "No module cache found"
```

## Cache Types and Expected Sizes

| Cache Type | Typical Size | Retention Strategy |
|------------|-------------|-------------------|
| Go modules (`~/go/pkg/mod`) | 100-500 MB | Keep latest per Go version |
| Go build cache (`~/.cache/go-build`) | 50-200 MB | Clean after 7 days |
| Docker buildx | 500 MB - 2 GB | Use `mode=min`, clean duplicates |
| Pre-commit hooks | 50-100 MB | Keep latest per config hash |
| Python pip cache | 20-50 MB | Keep latest per requirements hash |

## Best Practices

### 1. Cache Key Naming
- Include relevant version numbers (Go version, Node version, etc.)
- Use specific file hashes (go.mod, package.json) instead of wildcards
- Implement hierarchical restore keys for better hit rates

### 2. Regular Maintenance
- Monitor cache usage weekly
- Run aggressive cleanup when approaching 8 GB
- Review and update cache strategies quarterly

### 3. Emergency Procedures
If cache usage exceeds 9.5 GB:

1. **Immediate action**: Run emergency cleanup
   ```bash
   gh workflow run cache-cleanup.yml -f cleanup_strategy=emergency
   ```

2. **Manual cleanup**: Use the cleanup script
   ```bash
   ./scripts/cleanup-caches.sh stuartshay/gcp-automation-api 6 aggressive
   ```

3. **Temporary measures**: Disable caching in workflows temporarily
   ```yaml
   # Comment out cache steps temporarily
   # - name: Cache Go modules
   #   uses: actions/cache@v4
   ```

## Monitoring and Alerts

### GitHub Actions Cache API
Monitor cache usage programmatically:

```bash
# Get total cache size
gh api repos/stuartshay/gcp-automation-api/actions/caches | \
jq '[.actions_caches[].size_in_bytes] | add / 1024 / 1024 / 1024'
```

### Workflow Outputs
Check workflow logs for cache size information in the "Check cache sizes" steps.

## Troubleshooting

### Common Issues:

1. **Cache not found**: Normal for first runs or after cleanup
2. **Build slower after cleanup**: Expected temporarily, caches will rebuild
3. **GitHub API rate limits**: Wait and retry, or use personal access token

### Performance Impact:
- **Conservative cleanup**: Minimal impact, only removes very old caches
- **Aggressive cleanup**: Temporary slowdown (1-2 builds) while caches rebuild
- **Emergency cleanup**: Significant temporary slowdown, all caches need rebuilding

## Future Improvements

1. **Predictive cleanup**: Analyze cache usage patterns to predict optimal cleanup timing
2. **Selective caching**: Implement conditional caching based on change detection
3. **Cache warming**: Pre-populate caches for common scenarios
4. **Multi-stage cleanup**: Implement gradual cleanup based on cache age and usage frequency

## References

- [GitHub Actions Cache Documentation](https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows)
- [GitHub CLI Cache Commands](https://cli.github.com/manual/gh_api)
- [Docker Buildx Cache Documentation](https://docs.docker.com/build/cache/)
