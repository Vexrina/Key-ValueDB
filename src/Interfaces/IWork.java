package Interfaces;

public interface IWork<K, V> {
    void add(K key, V value, int time);
    void delete(K key);

    void put(K key, V value);
    V get(K key);
    int size();
}
