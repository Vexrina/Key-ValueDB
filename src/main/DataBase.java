package main;

import Interfaces.IDataBase;
import Interfaces.ITable;
import java.util.HashMap;
import java.util.Map;

public class DataBase implements IDataBase {
    private Map<String , ITable> workMap;

    public DataBase() {
        this.workMap = new HashMap<String, ITable>();
    }

    @Override
    public void add(String key, ITable value) {
        this.workMap.put(key, value);
    }

    @Override
    public ITable get(String key) {
        return this.workMap.get(key);
    }

    @Override
    public void remove(String key) {
        this.workMap.remove(key);
    }

    @Override
    public void put(String key, ITable value) {
        this.workMap.replace(key, value);
    }
}